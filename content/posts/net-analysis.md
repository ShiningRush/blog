+++
title = "网络分析"
date = "2021-02-12T15:09:23+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["网络分析"]
keywords = ["网络分析"]
description = ""
showFullContent = false
+++

# 网络分析

在 `linux` 上我们可以使用 `tcpdump` 来分析流量包，`wireshark` 分析包，`strace` 查看进程调用。
`fiddler` 用于分析 http 协议，`wireshark` 用于分析 tcp/udp 的网络封包

## https 包抓取
网络流量的这些工具为了能够解析 https 的包，通常都是自己签发证书，然后让系统信任自己证书，以作为中间人去转发、解析 https 流量。
如果仅仅只是为了代理 https 流量而不解析，可以使用 http 的隧道协议，使用 `CONNECT` method 去连接代理服务器，然后代理服务自动转发握手请求，相当于客户端直接与目标服务连接。

## NAT(Net Address Translation) 
NAT 常用于虚拟化技术中，又分为三类
- 静态NAT：此类NAT在本地和全局地址之间做一到一的永久映射。须注意静态NAT要求用户对每一台主机都有一个真实的Internet IP地址。
- 动态NAT：允许用户将一个未登记的IP地址映射到一个登记的IP地址池中的一个。采用动态分配的方法将外部合法地址映射到内部网络，无需像静态NAT那样，通过对路由器进行静态配置来将内部地址映射到外部地址，但是必须有足够的真正的IP地址来进行收发包。
- 地址端口NAT（NAPT）：最为流行的NAT配置类型。通过多个源端口，将多个未登记的IP地址映射到一个合法IP地址（多到一）。使用PAT能够使上千个用户仅使用一个全局IP地址连接到Internet。NAPT 又根据转换源头或者目标不同又分为 SNAT 和 DNAT，通常都会一起使用。实现 NAPT 的常见手段是设备维护一张转换表，里面存储了由 源地址：端口 到 转换后地址：端口 的映射关系

## tcpdump用法小结
### 常用选项
- `-i any` 可以抓取所有网卡流量
- `src/dst host hostname` 抓取 来自(src) or 发往(dst)，且 host 为 hostname 的流量
- `-s 0` 默认只抓取 68 字节，指定 0 可以抓取完整数据包
- `-X` 以十六进制的方式查看数据包细节
- `port 8080` 抓取端口 8080 的流量
- `tcp/udp port 8080` 指定 8080 端口的 tcp/udp 流量

### 输出分析
一般格式为
> TIME SRC > DST: Flags [.], data-seq, ack-seq, win, options, length

- TIME: 时间戳
- SRC: 源头
- DST: 目标地址
- Flags: 表示 TCP 数据包的标志位，F-FIN, S-SYN, .-ACK, S.-SYN+ACK, P.Push+ACK, R-RST 连接重置, 
- data-seq: 数据序列号，包括起始以及结束
- ack-seq: 已接受的序列号表示期望从其之后开始，
- win: 滑动窗口的缓存大小
