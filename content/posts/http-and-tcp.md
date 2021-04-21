+++
title = "http和tcp连接小记"
date = "2021-02-12T14:41:42+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["http", "tcp"]
keywords = ["http", "tcp"]
description = ""
showFullContent = false
+++
# http和tcp连接小记
## TCP 短连接与长连接
短连接与长连接的关键区分在于，`是否复用TCP连接`，短连接每次使用完毕都会释放连接，长连接则会将连接放入池中，等待下次使用。
所以实现长连接一般都会有几个关键参数：
- `最大连接数`: 其实就是连接池中保存连接数上限
- `过期时间`: 当请求数下降时，我们没有必要再维护这么多连接，所以我们需要给每条连接设置过期时间，当过期后及时从池中移除

在 `http1.0` 协议中，需要指定 `Connection: keep-alive` 来启用 TCP 长连接，`http1.1` 中已经成为默认值，无需指定，但是很多浏览器依旧会发送这个字段来兼容一些遗留的 web 服务。

*这里可能需要区分另一个用于 watch 结果的实现方式，我们知道一般有三种方式: 短轮询, 长轮询, websockt，长轮询当中会把请求 hold 住一段时间，让其也变成了 http 协议上的长连接，但是并不是我们 讨论的 tcp 协议上的长连接*

## 长连接消息的传输
传输编码：`Transfer-Encoding`， 内容编码: `Content-Encoding`
在长连接当中，由于连接不会断开，所以客户端无法知道消息是否已发送完毕，所以出现了 `Content-Length` 字段来告诉浏览器消息已传输完毕。但是很多情况下要计算出完整的响应体都是一件不容易的事情，为了计算出这个字段，我们需要把内容全部加载到内容中，这种实现非常不友好。所以出现了 `Transfer-Encoding`，一般为 `chunked` 代表分块传输。它不但解决动态内容长度的问题，同时还可以配合内容编码 (一般为 `gzip` or `deflate`) 进行边压缩边发送而无需预先知道压缩之前的内容大小。

### 分块格式
如果一个HTTP消息（包括客户端发送的请求消息或服务器返回的应答消息）的Transfer-Encoding消息头的值为chunked，那么，消息体由数量未定的块组成，并以最后一个大小为0的块为结束。
每一个非空的块都以该块包含数据的字节数（字节数以十六进制表示）开始，跟随一个CRLF （回车及换行），然后是数据本身，最后块CRLF结束。在一些实现中，块大小和CRLF之间填充有白空格（0x20）。
最后一块是单行，由块大小（0），一些可选的填充白空格，以及CRLF。最后一块不再包含任何数据，但是可以发送可选的尾部，包括消息头字段。
消息最后以CRLF结尾，
例子：
```
25
This is the data in the first chunk

1C
and this is the second one

3
con

8
sequence

0
```

更详细内容参考：[分块传输编码](https://zh.wikipedia.org/wiki/%E5%88%86%E5%9D%97%E4%BC%A0%E8%BE%93%E7%BC%96%E7%A0%81), [HTTP 协议中的 Transfer-Encoding](https://imququ.com/post/transfer-encoding-header-in-http.html)

很多语言都封装了完善的库来帮助我们完成 http 请求，导致我们很多时候都不知道 http 传输时的一些细节。实际上，这相当重要，会影响到 http 传输的性能。比如：
- 给出更大的 buffer 空间，你就会赢得更快的传输速度 (不超过带宽)，但是牺牲更多的空间。
- 如果 buffer 小于 chunk 的大小，那么会导致单个块传输效率更低甚至在某些库下会发生错误。

## TIME_WAIT 状态
`TIME_WAIT`状态是在四次挥手后，主动断开连接的一方会出现的状态，通常等待两个MSL时间 (*http规范为2分钟，但现在多数实现为30秒，所以一般等待一分钟*)会进入最后的`CLOSED`状态。
被动断开的一方会直接进入`CLOSED`状态。

`Nginx`中设置长连接可以有效避免 `TIME_WAIT` 产生；之前很奇怪为什么连接上游服务时，`Nginx`却会产生`TIME_WAIT`的连接，后来仔细想了想，一般场景中，无论调用方还是请求方，对连接都会有`Close`的调用，防止内存泄漏，所以谁先关闭都是有可能的。

`Nginx`设置长连接可以参考：[Nginx Upsteam](http://nginx.org/en/docs/http/ngx_http_upstream_module.html#keepalive) 以及 [支持keep alive长连接](https://skyao.gitbooks.io/learning-nginx/content/documentation/keep_alive.html)

值得注意的是，`Nginx`的`keepalive_requests`要和`keepalive`一起使用才会有效果。

`keepalive_requests`指当长连接被使用多少次之后就会释放该连接，按照官方文档说法，周期性地释放连接是必要的，这样才可以释放每个连接所申请的内存。
`keepalive`指保持多少空闲连接在缓存中。

## SYN_RECV 状态
当 Client 发送了 SYN 包，SEVER 端回复 ACK 则会进入该状态。
目前，Linux下默认会进行5次重发SYN-ACK包，重试的间隔时间从1s开始，下次的重试间隔时间是前一次的双倍，5次的重试时间间隔为1s, 2s, 4s, 8s, 16s，总共31s，第5次发出后还要等32s都知道第5次也超时了，所以，总共需要 1s + 2s + 4s+ 8s+ 16s + 32s = 63s，TCP才会把断开这个连接。由于，SYN超时需要63秒，那么就给攻击者一个攻击服务器的机会，攻击者在短时间内发送大量的SYN包给Server(俗称 SYN flood 攻击)，用于耗尽Server的SYN队列。对于应对SYN 过多的问题，linux提供了几个TCP参数：tcp_syncookies、tcp_synack_retries、tcp_max_syn_backlog、tcp_abort_on_overflow 来调整应对。

## Socket连接
可以由一个五元组来定义`[ 源IP, 源端口, 目的IP, 目的端口, 类型：TCP or UDP ]`，每个连接在 `Unix` 中会占用一个文件描述符(`FD`)

关于`TIME_WAIT`更详细的内容可以参考：[系统调优你所不知道的TIME_WAIT和CLOSE_WAIT](https://zhuanlan.zhihu.com/p/40013724)

## LVS(IPVS) 的长连接问题
Linux 对于长连接有自己的保活机制，在一个 TCP 连接一段时间不活动后则尝试发送探测包，如果探测失败则移除连接
查看 Linux 保活时间
```
sysctl -a | grep keep
net.ipv4.tcp_keepalive_intvl = 75  # 长连接探测包的间隔(s)
net.ipv4.tcp_keepalive_probes = 9  # 判断为失败的次数
net.ipv4.tcp_keepalive_time = 7200 # 空闲多少秒之后开始进行探测
```
值得一提的是KeepAlive并不是 TCP 规范的一部分，但是基本上操作系统都实现了它。但是光操作系统层面的设置其实对于应用层来说是不够的：
- 这几个配置都在内核参数中，配置繁琐
- 系统层面的可用其实不代表物理层面的可用，有的 ESTABLISHED 连接其实已经不可用了，但是应用层感知不到

为了解决这些问题，很多框架和库都在应用实现了自己的保活机制( KeepAlive )，本质上与系统的机制类似，也是通过心跳包来探测连接是否可用。

当使用 LVS 作为四层代理时，如果Linux 的TCP 保活时间高于 LVS 保活时间，那么在LVS主动断开连接后，Linux却不知道连接已断开，这时候分几种情况：
- 去 read 该连接会得到一个 `Connection reset by peer` or `EOF` 错误，这个取决于当前连接状态，暂时没找到靠谱的答案，
- write 会产生 RST 消息，如果继续 write 会产生 `broken pipe` 的错误。

查看 lvs 保活时间
```
ipvsadm --list --timeout
```
查看 lvs 连接
```
ipvsadm -lnc
```

这里有一个关于 Mysql 的 [典型超时问题](https://github.com/go-sql-driver/mysql/issues/257)
推测是 Mysqld 主动切断客户端连接，但是客户端却没有感知到。

## TCP连接超时问题
通常多数 TCP 连接都会依赖 OS 的默认超时时间，那么 OS 的默认超时时间时多少呢，使用以下命令能查看 TCP 重试相关的次数，而每次重试都是上一次重试时间的一倍，即遵从`1s 2s 4s` 这样的规律。
这种间隔重试方法被称为 `指数退避算法`
```
sysctl -a | grep retries
```

## http 协议: 100-Continue
当 client 需要在 body 放入数据时，有时候数据较大，需要检查服务端是否接受，可以采用使用该协议。
客户端可以在请求头部携带头部 `Expect: 100-continue`，curl 命令在 body 大于 1024 字节时会自动携带该头部。

注意以下几点：
- 由于 Server 端不一定会正确处理 100 协议，因此 client 应该在指定的 timeout 之后立即发送body
- Server 接受到 Expect 请求应该先响应 StatusCode 100 再继续读取请求体

## http samesite 限制
在chrome 51 之后为了限制 CSRF 攻击，限制跨站 cookie 的访问，如果目标站点设置的 Cookie 没有指定 `SameSite` 会被认为成 `Lax`
`SameSite` 分为三个等级：
- None: 不限制
- Lax: Cookies允许与顶级导航一起发送，并将与第三方网站发起的GET请求一起发送。这是浏览器中的默认值。
- Strict: Cookies只会在第一方上下文中发送，不会与第三方网站发起的请求一起发送。

## 网络分层
- OSI参考模型中分为了七层：应用、表示、会话、传输、网络、数据链路、物理，TCP/IP 简化为了四层：应用、传输、网络、物理
- OSI参考模型一直没有被广泛采用的原因就是因为它诞生于实验室，缺少真实的实践而TCP/IP 则是从实践中沉淀而来
- OSI参考模型一般用于在讨论使用

## IP 地址和 MAC 地址的作用
IP 地址用于网络层的寻址，MAC 地址用于数据链路层的寻址。
一个数据包由源地址发往目标地址的时都会经由好几个路由设备，在这个过程中，IP是不变的，但是MAC地址却不断在变化。
而每经由一层链路层的数据包转发，都会被称为一跳( Hop )

在最开始互联网定义几类标准的IP地址（A，B，C，D），后来发现这些满足不了具体的使用场景，于是诞生了CIDR（Classless In-Domain Routing），通过子网掩码来区分网路标识和主机标识，如 9.128.100.0/24，表示前24位用于网络标识，而剩下的作为主机标识。
