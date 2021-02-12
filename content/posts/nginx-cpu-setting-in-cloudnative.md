+++
title = "Nginx在云上环境的CPU最佳实践"
date = "2021-02-12T15:12:53+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["nginx", "openresty"]
keywords = ["nginx", "openresty"]
description = ""
showFullContent = false
+++

# Nginx在云上环境的CPU最佳实践
## 目录
- [背景](#背景)
- [Nginx的CPU相关设置](#Nginx的CPU相关设置)
- [k8s的CPU策略](#k8s的CPU策略)
- [方案对比](#方案对比)
- [结论](#结论)

### 背景
最近有用户反馈 [ApacheAPISIX](https://github.com/apache/apisix) 在云原生环境下存在几个问题：
- 无法获取容器准确核数
- 配置了多核的情况下无法完全利用

这里先简单介绍项目相关的情况，APISIX 是一个基于 Openresty 的开源网关。而 Openresty 其实就是 Nginx + LuaJIT，那么我们要调查问题其实跟 Nginx 是脱不开关系的。

### Nginx的CPU相关设置
首先看看第一个问题：**无法获取容器准确核数**，这个问题的起因是因为在 Nginx 配置中使用了
```
worker_processes auto;
```
`auto` 意味 Nginx 会自动获取 CPU 核数，然后根据核数创建 worker。不幸的是，在容器当中它获取到的是 **母机的核数** 导致 Nginx 会在容器中创建数十个甚至上百个 worker，多个 worker 间的资源竞争和上下文切换都会降低它的性能。
为了核验 Nginx 是否真的获取的是母机核数，我翻了下 Nginx 相关的代码，截取核心片段如下：
src/os/unix/ngx_posix_init.c
```c
#include <ngx_config.h>
#include <ngx_core.h>
#include <nginx.h>


ngx_int_t   ngx_ncpu;
ngx_int_t   ngx_max_sockets;
ngx_uint_t  ngx_inherited_nonblocking;

...

#if (NGX_HAVE_SC_NPROCESSORS_ONLN)
    if (ngx_ncpu == 0) {
        ngx_ncpu = sysconf(_SC_NPROCESSORS_ONLN);
```
查看 [sysconf](https://man7.org/linux/man-pages/man3/sysconf.3.html) 的文档发现底层调用的 [ get_nprocs_conf(3)](https://man7.org/linux/man-pages/man3/get_nprocs_conf.3.html), 继续查看它的[源码](https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/getsysstats.c.html#__get_nprocs_conf)，核心片段如下：
```c
/* On some architectures it is possible to distinguish between configured
   and active cpus.  */
int
__get_nprocs_conf (void)
{
  /* XXX Here will come a test for the new system call.  */
  /* Try to use the sysfs filesystem.  It has actual information about
     online processors.  */
  DIR *dir = __opendir ("/sys/devices/system/cpu");
  if (dir != NULL)
  ...
```
注意这个路径 `/sys/devices/system/cpu`，随便进入到一个容器中，ls 一下它你会发现它是母机的CPU信息，类似下面：
![cpulist](/image/nginx-cpu-setting-in-cloudnative/cpulist.png)

OK，第一个问题至此已确认完毕，我们先看看第二个问题再讨论解决方案，因为看上去它们应该是关联的。

### k8s的CPU策略
关于第二个问题：**配置了多核的情况下无法完全利用**，直觉这个问题跟另一个 Nginx 的配置参数有关：`worker_cpu_affinity`，它可以指定 Nginx绑定到几号核。手动指定的场景我们直接跳过，通常我们都是使用 `auto` 参数。看下当设置 auto 时，Nginx 会怎么绑定CPU，核心代码片段如下：
src/core/nginx.c
```c
    if (ngx_strcmp(value[1].data, "auto") == 0) {

        if (cf->args->nelts > 3) {
            ngx_conf_log_error(NGX_LOG_EMERG, cf, 0,
                               "invalid number of arguments in "
                               "\"worker_cpu_affinity\" directive");
            return NGX_CONF_ERROR;
        }

        ccf->cpu_affinity_auto = 1;

        CPU_ZERO(&mask[0]);
        for (i = 0; i < (ngx_uint_t) ngx_min(ngx_ncpu, CPU_SETSIZE); i++) {
            CPU_SET(i, &mask[0]);
        }

        n = 2;

    } else {
        n = 1;
    }
```
可以看到CPU的绑核策略是顺序从低位到高位，这样做在普通的物理机本来没什么问题，但是在 k8s 的环境下就不行了，原因有两个：
- 绑核动作需要特权执行，通常 POD 是没有权限的
- 在于 k8s 在`static` 策略下本来就会对 limit 为整数的 `Guaranteed` POD进行绑核处理，可以参考 [控制节点上的 CPU 管理策略](https://kubernetes.io/zh/docs/tasks/administer-cluster/cpu-management-policies/) 。

所以云上应用都不建议再去应用进行绑核操作。

回过头来说 Nginx 自动从低位CPU绑到高位的这个操作，没有特权的情况会怎样？上面的代码片段我们看到它使用了 `CPU_SET` 这个系统调用，相关的方法签名如下：
```
void CPU_SET(int cpu, cpu_set_t *set);
```
意味着无论绑核成功或失败，**程序都得不到响应**。为了验证这个结论，我们创建一个 Nginx 应用( 1c1g )，然后在容器执行以下命令查看绑定的核：
```
sh-4.2# cat /sys/fs/cgroup/cpuset/cpuset.cpus
45
```
可以看到绑定到了第45号核(由0开始)，在母机上执行 `htop` 可以看到这里的第46号核(由1开始)，完全没有使用率：
![htop-before](/image/nginx-cpu-setting-in-cloudnative/htop-before.png)

开始执行压测后：
![htop-after](/image/nginx-cpu-setting-in-cloudnative/htop-after.png)

很明显，Nginx 绑核并没有成功，容器依然绑定在原来的CPU上。
通常来说，没有特殊原因都不建议云上应用再去执行绑核操作，保持不变即可。

### 方案对比
开始我以为多核无法利用的情况，是容器绑定的核与应用绑定的核只有小部分重叠，所以才导致无法有效利用。但现在看来，没有 **特权** 的 Nginx 甚至连绑核都做不到，那么我们要继续考虑其他可能的问题。
列举以下情况：
| 情况 | QPS | CPU使用率 |
| --- | --- | --- |
| 1c2g-1workers | 10959.22 | 100% | 
| 2c4g-1workers | 11845.91 | 100% |
| 2c4g-2workers | 16975.04 | 200% |
| 1c2g-1workers * 2 | 22492.83 | 200% |
| 4c8g-2workers | 20506.10 | 200% |
| 4c8g-4workers | 31012.40 | 400% |
| 1c2g-1workers * 4 | 51720.11 | 400% |
| 1.001c2g-1workers * 4 | 47501.16 | 401% |

这里我们可以看到使用多个单核的容器会比使用多核的单个容器具有更高的吞吐量，可能的原因有两点：
- 多个worker间存在资源竞争
- 单核容器由于可用核只有一个，所以相当于进行了绑核操作

为了排除第二个原因的影响，我补充图表中的最后一个用例(1.001c2g-1workers * 4)，这个用例中由于cpu不是整数核，所以没有分配独占cpu。可以发现它的吞吐量下降了8%左右，证明绑核还是有效果的。

### 结论
通过以上的实验，我们可以得出几个关键点：
- worker 少于 cpu 数量时无法充分利用 cpu
- 云上的 Nginx 最好使用 **单核多容器** 的部署模式，这样可以充分利用k8s的cpu策略进行绑核，单worker也是Nginx的推荐设置，如果要设置多个 worker，那么不要使用 `auto` 参数，获取到的将会是母机核数
- 如果要使用多核的情况下，尽量不要使用过大的cpu数量，推荐最多8个，否则 worker 数量过多会造成大量资源浪费在处理竞争上
- 云上不要使用 `worker_cpu_affinity=auto`，因为除了需要特权外，顺位绑核的操作不一定能绑定到 k8s 分配的独占核，极端情况下还会导致不可用，所以在云上环境多结合 k8s 的 CPU策略 来使用才是最佳实践