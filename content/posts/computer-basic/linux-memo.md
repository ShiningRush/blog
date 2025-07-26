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
- 查看cpu信息
```
top/htop/dstat
```
- 查看网卡信息
```
watch -n 1 "/sbin/ifconfig eth0 | grep bytes"
```

## 常见信息
linux 系统下的日志基本都在 `/var/log/` 目录下，有系统写的，也有服务写。最常用的就是 `/var/log/messages`(系统日志), `/var/log/dmesg`(内核日志)，日志目录下的各类日志都是很多通过 `syslog` 的进程来产生。
- 内核日志: `/var/log/dmesg_all`, `/var/log/dmesg`
- 系统日志： `/var/log/messages`，部分系统 会写入 `/var/adm/messages`

你也可以使用 `dmesg` 命令来读取内核 ring buffer 中的内容，它会在下次启动前输出到 `/var/log/dmesg` 中，因此 `dmesg` 命令更实用
参考：[Difference between output of dmesg and content of /var/log/dmesg?](https://unix.stackexchange.com/questions/191560/difference-between-output-of-dmesg-and-content-of-var-log-dmesg)

systemd 所维护的程序日志可以通过 `journalctl -u xxxx` 指令来查看，比如查看 kubelet 的日志：`journalctl -u kubelet`，日志根据配置
```bash
/etc/systemd/journald.conf
/etc/systemd/journald.conf.d/*.conf
/run/systemd/journald.conf.d/*.conf
/usr/lib/systemd/journald.conf.d/*.conf
```
可能会存在不同地方，默认在 `/var/log/journal` or `/run/log/journal`，它并不会持久化，在下一次程序启动后就会丢失。

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


## 查看 unix socket 的peer
下面的命令可以查看所有的UDS
```bash
ss -a --unix
```

## 程序coredump
Core Dump 又叫核心转储。在程序运行过程中发生异常时，将其内存数据保存到文件中，这个过程叫做 Core Dump。Core是指记忆体也就是现在的内存。
可以使用
```shell
# 查看限制
ulimit -a

# 打开coredump限制
ulimit -c unlimited

# 查看路径和格式, 不存在的话默认写到程序所在目ls
# 如果开头为 | 则把剩余部分视为一个程序，将coredump文件作为标准输入调用
# 参考ttps://stackoverflow.com/questions/47765202/what-does-mean-in-file-proc-sys-kernel-core-pattern
cat /proc/sys/kernel/core_pattern

# 控制core文件名是否包含pid，默认为0
cat /proc/sys/kernel/core_uses_pid

# 分析coredump文件，在gdp中使用 where or bt 来查看崩溃时的信息
#2  `p $_siginfo` 可以查看具体的 signal information
#3  `x/i $pc` 查看core执行的汇编指令
#4  `i r` 查看寄存器值
gdb -c [core_file] [bin] 

```

注意分析coredump文件时可能出现问号，这是由于编译时没有带上符号链接所导致(-g)。不过带上符号链接也有几个风险：
1. 体积会增大很多
2. 会泄漏源码（反编译）

因此如果程序运行在不安全的环境，最好是编译额外的map文件来配合分析调查。

## 调查命令
如果需要查询pipe对端的进程，使用以下命令：
```bash
(find /proc -type l | xargs ls -l | fgrep 'pipe:[20043922]') 2>/dev/null
```