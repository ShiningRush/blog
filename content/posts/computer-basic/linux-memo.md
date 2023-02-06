+++
title = "Linux系统小记"
date = "2021-02-12T14:47:09+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["linux"]
keywords = ["linux"]
description = ""
showFullContent = false
+++

# Linux系统小记
## 修改环境变量
- 可以通过 `export var=value` 在当前控制台修改，不会永久生效
- 可以通过 以下文件 修改 ，添加 `export var=value`
```
============
/etc/profile
============
此文件为系统的每个用户设置环境信息,当用户第一次登录时,该文件被执行.
并从/etc/profile.d目录的配置文件中搜集shell的设置.

===========
/etc/bashrc
===========
为每一个运行bash shell的用户执行此文件.当bash shell被打开时,该文件被读取.

===============
~/.bash_profile
===============
每个用户都可使用该文件输入专用于自己使用的shell信息,当用户登录时,该
文件仅仅执行一次!默认情况下,他设置一些环境变量,执行用户的.bashrc文件.

=========
~/.bashrc
=========
该文件包含专用于你的bash shell的bash信息,当登录时以及每次打开新的shell时,该文件被读取.

==========
~/.profile
==========
在Debian中使用.profile文件代 替.bash_profile文件
.profile(由Bourne Shell和Korn Shell使用)和.login(由C Shell使用)两个文件是.bash_profile的同义词，目的是为了兼容其它Shell。在Debian中使用.profile文件代 替.bash_profile文件。

==============
~/.bash_logout
==============当每次退出系统(退出bash shell)时,执行该文件
```

## 常用命令
- 查找文件是否包含字符串
```
grep -rn "string" /path/to/search
```
- 查找文件名
```
find /path/to/search -name "string"
```
- 查看磁盘信息
```
df -h
```
- 查看文件夹大小（前者仅仅查看目录，后者查看目录以及所有子项）
```
du -sh /path/to/look
du -h /patt/to/look
```
- 查看 tcp 网络
```
netstat -atlnp
```

## 查看相关信息
-查看CPU型号
```
cat /proc/cpuinfo | grep name | cut -f2 -d: | uniq -c
```
- 查看物理CPU的个数
```
cat /proc/cpuinfo |grep "physical id"|sort |uniq|wc -l
```
- 查看逻辑CPU的个数
```
cat /proc/cpuinfo |grep "processor"|wc -l
```
- 查看每个CPU中core的个数(即核数)
```
cat /proc/cpuinfo |grep "cores"|uniq
```
- 查看操作系统内核信息
```
uname -a
```
- 查看CPU运行在32bit还是64bit模式
```
getconf LONG_BIT
```
- 查看内存总量
```
grep MemTotal /proc/meminfo  
```
- 查看空闲内存总量
```
grep MemFree /proc/meminfo
```

## 常见信息
系统日志一般位于 `/var/log` 中，
- 内核日志: `/var/log/dmesg_all`, `/var/log/dmesg`
- 系统日志： `/var/log/messages`，部分系统 会写入 `/var/adm/messages`

## 虚拟内存
进程在操作内存时并不是直接操作物理内存，而是操作系统给进程模拟的 `虚拟内存`，使用这种技术有几个优点：
- 避免竞争：进程不需要考虑与其他进程的竞争问题
- 内存完整性：进程看到的内存都是完整而连续的
- 安全性：操作系统可以通过内存页来控制内存的可访问性
- 数据共享：可以让不同虚拟内存映射到相同物理内存上来实现内存共享

参考[linux内存管理](https://zhuanlan.zhihu.com/p/149581303)

Tips: 这和 Windows上可以手动分配的 `虚拟内存`功能不同， windows 的 虚拟内存更像是 SWAP分区

### 内存页
我们知道数据存储在内存上时最终都换转换为字节存储，它是内存存储的最小单元，在 32 位 系统下能支撑的内存大小为 `2^32Byte = 4GB`。
如果虚拟内存在为进程映射时为每一字节都建立映射关系，在 32 位系统下至少需要 `(4 + 4) * 2^32 Byte = 32GB` 的空间来存储映射数据，显然是不合理，所以提出了 `内存页` 来做映射关系。
Linux 在分配物理内存时，以页位单位分配，寻址时通过 页索引 + 偏移量 来映射，通常内存页大小在 4K (i386体系)

### SWAP分区
在物理内存完全使用完毕后，系统还可以通过将部分内存页（这部分数据来自一些很久没访问过的进程）转储到磁盘，让需要内存的进程先使用，等待进程需要再次访问这些数据时再交换回来。

## 日志分析
linux 系统下的日志基本都在 `/var/log/` 目录下，有系统写的，也有服务写。最常用的就是 `/var/log/messages`(系统日志), `/var/log/dmesg`(内核日志)，日志目录下的各类日志都是很多通过 `syslog` 的进程来产生。
systemd 所维护的程序日志可以通过 `journalctl -u xxxx` 指令来查看，比如查看 kubelet 的日志：`journalctl -u kubelet`，日志根据配置
```bash
/etc/systemd/journald.conf
/etc/systemd/journald.conf.d/*.conf
/run/systemd/journald.conf.d/*.conf
/usr/lib/systemd/journald.conf.d/*.conf
```
可能会存在不同地方，默认在 `/var/log/journal` or `/run/log/journal`，它并不会持久化，在下一次程序启动后就会丢失。

## 查看 unix socket 的peer
下面的命令可以查看所有的UDS
```bash
ss -a --unix
```

## 查看进程cgroup配置
通常查看整个系统的cgroup配置，可以直接查看 `/sys/fs/cgroup` 即可，但是我们可以自己创建cgroup节点，systemd也是采用类似的方式，因此要想查看一个进程级别的限制，我们可以通过如下方式：
1. 查看 `/proc/{pid}/cgroup`，我们可以得到如下的提示：
```
10:pids:/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399/system.slice/creativecloud.traffic.proxy.service
9:net_cls,net_prio:/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399
8:cpuset:/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399
7:memory:/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399
6:blkio:/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399
5:freezer:/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399
4:devices:/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399
3:perf_event:/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399
2:cpu,cpuacct:/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399
1:name=systemd:/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399/system.slice/creativecloud.traffic.proxy.service
0::/
```
2. 冒号后的内容即是对应controller的自定义节点，我们可以在
```
/sys/fs/cgroup/{controller}/{冒号后内容}
```
找到对应的控制节点，以cpuset为例，我们子系统的配置目录位于：
```
/sys/fs/cgroup/cpuset/kubepods/burstable/pod2f4f4d51-8f60-4a01-a9e2-3d35e54bf812/124f4cc0cef0835976ed1ad6ac88bae06c5b92b64c3f81f2a3c47de2a0d55399
```