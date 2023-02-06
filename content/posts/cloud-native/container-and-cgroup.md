+++
title = "cgroup与容器化"
date = "2022-12-18T11:33:56+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["", ""]
keywords = ["", ""]
description = ""
showFullContent = false
readingTime = false
+++
# cgroup与容器化
相信大家都知道 cgroup 是容器化的基础技术之一(这里我们指 `runC` 运行时，因为 `kata` 之类的容器技术使用的是虚拟化来进行隔离)，而现今很多互联网应用都在朝云原生进行改造，因此无论是否从事容器相关的工作，我们都可以了解下cgroup的工作原理，便于我们更了解自己程序的允许环境以及容器化技术是怎么保障服务间的隔离性的。

## cgroup是什么
cgroups 全称是 Linux Control Group，是内核的一个特性, 它运行研发以层次的结构来控制进程甚至是线程的资源用量（比如cpu、内存等），它的接口以一种伪文件系统(pseudo-filesystem)的形式暴露给用户使用(位于 `/sys/fs/cgroup` )。其中分组的层次关系由 `内核代码` 实现，而各个资源的限制与跟踪由其 `子系统(subsystem)` 来实现
> 详情可以参考 [man7.cgroup](https://man7.org/linux/man-pages/man7/cgroups.7.html)

## 概念解释
开始前我们先介绍下一些cgroup相关的术语：
- **cgroup**: 一个group，其中包含了受到相同资源限制的进程集合
- **subsystem**: 子系统会对应内核的一个组件，它负责调整cgroup中的进程行为，以使它们满足预期。linux中已经实现了各种各样的子系统，它们使得我们可以控制进程的各个资源项。有时子系统也被称为 `资源控制器(resources controller)` or `控制器`
- **hierarchy**: 上文提到过，cgroup 会按照层次结构来组成，而这个组成拓扑被称为 `hierarchy`, 而每一个层级都可以定义自己的属性而限制在它之下的子层级，因此更高层的属性定义不能超过它的后代总和，比如你的父层级cpu限制为2c，你的后代层级之和为4c。

## V1 与 V2版本
最初的cgroup 设计出现在 Linux2.6.24 中，随着时间发展，大量的控制器被增加，而逐渐导致控制器与hierarchy之间出现了巨大的不协调与不一致，详情可以参考内核文件：`Documentation/admin-guide/cgroup-v2.rst (or Documentation/cgroup-v2.txt in Linux 4.17 and earlier)`。
由于这些问题，linux从3.10开始实现一个正交的新实现(v2)来补救这些问题，直到 linux 4.5.x 才标记为稳定版本。文件 cgroup.sane_behavior，存在于v1 中，是这个过程的遗留物。该文件始终报告“0” 并且只是为了向后兼容而保留。
虽然 v2 版本是 v1 版本的一个替换，但是由于兼容性原因，它们都会持续存在于系统中，并且目前 v2 的控制器还只是 v1 版本的子集。当然，我们完全可以并行使用两者来完成我们的需求，唯一的限制是 v1 v2 的 hierarchy 不能存在相同的控制器（很容易理解，这是为了避免控制冲突。）

## cgroup v1
在 v1 版本下，每个控制器都会挂载一个独立的cgroupfs，提供各自的层次结构来管理系统上的进程。你还可以将多个子系统（甚至）所有子系统都挂载到一个目录下，而子 hierarchy 也会表现为目录下的子目录。以 `/user/joe/1.session` 为例，`1.session` 为cgroup，它是 `joe` 的子group，而 `joe` 为 `user` 的子group。

### Tasks(Threads) Vs Processes
在 v1 版本中，可以配置到 Task（在用户态，更常见的叫法是 Thread） 级别的资源控制，但是这也带来了一些问题，比如 `memory` 这个控制器，对于 Task 而言是没有意义的，因为线程之间都共享同一个内存空间。在 v2 版本从已经被移除，然后以受限的形式实现了类似的功能。(v2 - `ThreadMode`)

### 挂载 v1 控制器
开启 cgroup 能力需要在内核编译时携带 `CONFIG_CGROUP` 选项。为了使用 v1 控制器，它必须被挂载为 cgroup fs。通常位置是以 [tmpfs(5)](https://man7.org/linux/man-pages/man5/tmpfs.5.html) 挂载到 `/sys/fs/cgroup`。因此，一个可能的命令如下：
```
mount -t cgroup -o cpu none /sys/fs/cgroup/cpu
```
你也可以将两个控制器合挂到同一个 hierarchy 下，如：
```
 mount -t cgroup -o cpu,cpuacct none /sys/fs/cgroup/cpu,cpuacct
```

值得一提的是，大多数系统都已经实现了上面约定的目录，比如 `systemd` 就是如此。

### 卸载 v1 控制器
通常卸载控制器没什么特别的：
```
umount /sys/fs/cgroup/pids
```
但是要注意卸载前要确定 hierarchy 下已经没有了子group，否则只会导致目录不可见，而并不会自动卸载子group。

### 目前 v1 支持的控制器
这里就不搬运了，有兴趣的同学直接参考上面 man7 的页面。

### 创建 cgroup 并绑定进程
cgroup fs 默认初始化时会创建一个 root group `/`，所有的进程都归属这个它。如果你想要创建一个新的 group ，只需要在相应的 cgroup fs 创建一个新的目录即可，如：
```
mkdir /sys/fs/cgroup/cpu/cg1
```
这样即会创建出一个空的cgroup
而通过将 pid 写入 `cgroup.procs`，即可将进程与cgroup绑定。
```
echo $$ > /sys/fs/cgroup/cpu/cg1/cgroup.procs
```
> **Tips**
> 
> 注意一次只能写入一个pid

写入 `0` 到 `cgroup.procs` 则会将 **写入进程** 加入到对应的 cgroup。同时当进程加入 cgroup 时，其拥有的所有线程也会一起进入到cgroup。同时属于 `cgroup.procs` 下的进程可以通过该文件而获取到 pid list, 该列表不保证顺序，也可能会重复。
另外在一个 hierarchy 中，一个进程只能归属于 cgroup，因此当进程被加入某个 cgroup 时，也会从之前的grouo中移除。
你可以通过把 taskid 加入到 group 的 `tasks` 文件中来设置其归属的group。

### 删除 cgroup
确定其没有子group后，可以直接删除目录即可。

## cgroup v2
在 cgroup v2 中，所有挂载的控制器都会位于一个统一的 hierarchy。虽然不同的控制器可能同时处于 v1 于 v2 下，但是相同的控制器不能同时处于 v1 与 v2 下。
我们可以通过执行
```
mount | grep cgroup
```
来查看当前 cgroup 的挂载，一个开启了 v1 v2 系统如下所示：
```
tmpfs on /sys/fs/cgroup type tmpfs (ro,nosuid,nodev,noexec,mode=755,inode64)
cgroup2 on /sys/fs/cgroup/unified type cgroup2 (rw,nosuid,nodev,noexec,relatime,nsdelegate)
cgroup on /sys/fs/cgroup/systemd type cgroup (rw,nosuid,nodev,noexec,relatime,xattr,name=systemd)
cgroup on /sys/fs/cgroup/cpuset type cgroup (rw,nosuid,nodev,noexec,relatime,cpuset)
cgroup on /sys/fs/cgroup/cpu,cpuacct type cgroup (rw,nosuid,nodev,noexec,relatime,cpu,cpuacct)
cgroup on /sys/fs/cgroup/hugetlb type cgroup (rw,nosuid,nodev,noexec,relatime,hugetlb)
cgroup on /sys/fs/cgroup/blkio type cgroup (rw,nosuid,nodev,noexec,relatime,blkio)
cgroup on /sys/fs/cgroup/pids type cgroup (rw,nosuid,nodev,noexec,relatime,pids)
cgroup on /sys/fs/cgroup/rdma type cgroup (rw,nosuid,nodev,noexec,relatime,rdma)
cgroup on /sys/fs/cgroup/net_cls,net_prio type cgroup (rw,nosuid,nodev,noexec,relatime,net_cls,net_prio)
cgroup on /sys/fs/cgroup/freezer type cgroup (rw,nosuid,nodev,noexec,relatime,freezer)
cgroup on /sys/fs/cgroup/devices type cgroup (rw,nosuid,nodev,noexec,relatime,devices)
cgroup on /sys/fs/cgroup/memory type cgroup (rw,nosuid,nodev,noexec,relatime,memory)
cgroup on /sys/fs/cgroup/perf_event type cgroup (rw,nosuid,nodev,noexec,relatime,perf_event)
```
其中 `/sys/fs/cgroup/unified` 即是统一的 root hierarchy 

## 其他注意事项
通过 [fork2](https://man7.org/linux/man-pages/man2/fork.2.html) 创建的子进程会自动继承其父进程的的cgroup。如果是新启用的进程的cgroup传递使用 [execve(2)](https://man7.org/linux/man-pages/man2/execve.2.html)

### /proc 文件
#### /proc/cgroups
该文件包含了已编译到内核的控制器，如下：
```
subsys_name    hierarchy      num_cgroups    enabled
puset          4              1              1
cpu             8              1              1
cpuacct         8              1              1
blkio           6              1              1
memory          3              1              1
devices         10             84             1
freezer         7              1              1
net_cls         9              1              1
perf_event      5              1              1
net_prio        9              1              1
hugetlb         0              1              0
pids            2              1              1
```
从左往右，列的含义依次为：
1. 控制器的名称
2. 控制器挂载的 cgroup hierarchy 的唯一ID，如果 v1 的控制器挂载到了同一个 hierarchy，那么ID相同。以下情况值为0:
  1. 控制器没有挂载到 v1 的cgroups
  2. 控制器绑定到了 v2 的 unified hiararchy
  3. 控制器没有启用。
3. 在这个 hierarchy 下有多少cgroup
4. 是否禁用

#### /proc/[pid]/cgroup
这个文件描述了进程所属的 cgroup，显示的内容会因为 v1 v2的版本不同而不同。
对于该进程的每个 cgroup ，都会显示一条冒号分隔的信息：
> hierarchy-ID:controller-list:cgroup-path

举例：
> 5:cpuacct,cpu,cpuset:/daemons

从左往右，冒号分隔的三个字段含义为：
1. v1 版本中，该字段为 cgroup 所属的 root hierarchy(与 `/proc/cgroup` 中的id相同 )。对于v2版本，该字段为0。
2. v1 版本中，这个字段包含了一个逗号分隔的列表，用以描述绑定的控制器。对于v2版本，该字段为空。
3. 该字段显示了进程所属的cgroup在 hierarchy 中的路径，该路径是一个相对路径（相对于 hierarchy 的挂载点，比如上面的例子中，假设挂载点遵从默认路径 /sys/fs/cgrpup/cpu,cpuacct，那么这个cgroup的位置为 /sys/fs/cgrpup/cpu,cpuacct/daemons）

### /sys/kernel/cgroup 文件
#### /sys/kernel/cgroup/delegate
该文件暴露了 v2 版本中哪些内容是可以被委派的，未来有可能会改变，目前为：
```
cgroup.procs
cgroup.subtree_control
cgroup.threads
```

### /sys/kernel/cgroup/features
随着发展，v2 的功能可能会增加或者改变，因此这里提供了一个用户态的方式来发现哪些功能被启用了，目前内容如下：
```
nsdelegate
memory_localevents
```
