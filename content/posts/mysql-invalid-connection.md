+++
title = "mysql invalid connection 错误"
date = "2021-12-31T15:03:15+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["mysql"]
keywords = ["invalid connection"]
description = ""
showFullContent = false
+++

# mysql invalid connection 错误
最近团队有同学遇到了 mysql 会偶发"invalid connection" 错误，在之前从没遇到过，因此稍微调查一下。

## 原因
查看 [github.com/go-sql-driver/mysql@1.6.0](github.com/go-sql-driver/mysql) 的源码可以看出来会返回这个错误的地方不是很多，大致可以归类以下几个case：
- TCP连接已关闭：这里的连接关闭是指被动关闭，通常对数据库的连接池都由 go 的 `sql` 包来管理，无论是 `MaxIdleTime`, `MaxLifetime` 都是它来管理生命周期，如果取出了已过期的连接，它会重新再去取，然后关闭掉已过期的连接，这是主动关闭，源码参考 [sql](https://github.com/golang/go/blob/master/src/database/sql/sql.go#L1289)。被动的关闭，指从server端发起的，比如连接空闲时间 > `wait_timeout`，再比如 mysql or mysql proxy重启或者故障，这些都导致 `invalid connection`。从代码来看，导致这个错误的另一种情况是对一个已经关闭的`mysqlConn` 进行操作，这个case我想不到什么情况下会出现，因为逻辑上来说当这个对象被调用 `Close` 后，将没有机会被调用，不知道要什么情况会被复用到。
- 读写异常: 查看 `readPacket` 和 `writePacket` 代码可以发现，如果在读写过程发生错误，且 client 端没有 cancel context，那么会返回 `invalid connection`，具体的原因可以在标准输出里查看，mysql会输出 `[mysql]` 前缀的日志，最常见的是读写超时。
- 当开启插值参数后，并发使用连接：正常 mysql driver走的 prepare + exec 的过程，会产生多一次的 rtt，有些情况我们可以通过dsn参数 `interpolateParams=true` 来开启插值（占位符替换），而减少一次rtt。但是插值的过程会使用 conn 上的buffer，如果并发访问导致 buffer 繁忙则会返回该错误。那么什么时候会在并发下使用同一条连接呢？答案是你当你使用了一些 session性质的功能，比如 transaction or user-level lock。

插值相关的参考资料：[database/sql 一点深入理解](https://michaelyou.github.io/2018/03/30/database-sql-%E4%B8%80%E7%82%B9%E6%B7%B1%E5%85%A5%E7%90%86%E8%A7%A3/)

## 解决方案
首先要确定自己是由于何种原因导致的错误，上面的原因当中，除了第一种不会有日志外，其他都可以在标准输出种找到 mysql 的日志，如果没办法核获取日志，可以按照以下思路排查：
- 如果是连接关闭，那么检查 mysql or mysql 的proxy是否异常重启过，mysql的wait_timeout是否配置得过小，或者等于client连接的Lifetime。
- 如果是读写异常，原因可能很多种，建议找到日志再去排查。但是最常见还是读写超时，可以看看读写超时的配置时间是否不合理。
- 如果你没有在dsn开启 `interpolateParams=true`，那么第三种情况是不存在的；如果开启了，注意不要在session操作中使用并发。

原因当中的第一种和第二种其实有关联的，比如读写过程中，连接突然被关闭，也会得到一个 invalid connectio。

## 小结
超时和并发的原因两个原因比较直接而明显，这里不多解释了，值得一提的，TCP连接被动关闭后，应用层是没有感知的，导致使用了一个无效连接写了一堆东西之后会得到 `broken pipe`的错误(go 是这样，tcp来看其实是得到了一个RST)，进而引发错误。
针对这种情况，mysql 的 driver 在 2019年（虽然问题很早就已经暴露了） 添加了 tcp 的存活检测 [#934](https://github.com/go-sql-driver/mysql/pull/934)，如果发现tcp连接已经被关闭后，会返回一个 `driver.BadConn`，`sql` 包会针对这个错误进行重试。
放大来看，其实所有go的web框架在长连接上都面临这种问题，存活检查是不得不考虑的，无论是 `before-use` 还是 `heartbeat`，否则TCP连接预期外的被动关闭时就只能产生错误。