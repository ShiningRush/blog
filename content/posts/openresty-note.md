+++
title = "Openresty 笔记"
date = "2021-02-12T15:10:35+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["openresty"]
keywords = ["openresty"]
description = ""
showFullContent = false
+++

# Openresty 笔记
## 简介
Openresty 是一个国人将LuaJIT嵌入Nginx进程进而可以使用Nginx来进行开发高性能的Web框架。
入门的简介可以参考这个文档，[OpenResty 不完全指南](https://juejin.im/entry/5ba3abd65188255c8a05f69c)

## 注意点
- Openresty 由于是在每个 nginx worker 都运行了一个 luajit 所以它仍然是一个单线程模型所以不存在 `并行(parall)` 问题，但是用于处理请求的是 lua 的 协程`coroutine`，那么这里可能存在 `并发(parallelism)` 问题，即如果 协程A 因为某种原因挂起后yield(IO操作/sleep)，此时 协程B 将会被运行，同时共享 `模块变量` 以及 `全局变量`。openresty 这种非抢占式的协程调度带来了优点，也有缺点，优点就是不存在多个协程同时进行导致的并发读写问题，缺点在于如果某个协程进行一些长时间的调用，比如计算或者系统调用等，那么整个进程都会被卡主，导致吞吐量大大下降。（go 在 1.14 之后实现了基于信号的抢占式协程调度，由 `sysmon` 线程来执行检测，每 20us 检测周期，如果发现P的running状态超过 10ms or syscall 超过 20us，那么直接发送信号给 M，让其放弃正在执行的 G，开始执行其他的 G）
- 使用 Openresty 可以调用 c 的代码，底层是利用了 linux 的 `dlsym` 和 `dlopen` 利用后者加载链接库，前者加载函数符号。由于 LuaJIT 与 C 进行交互，所以不会识别 c++ 的机制，比如析构函数，这样一来你就没法使用 c++ 推荐的 `RAII` 机制。你有两个方式来管理内存：1、c way，在使用完后释放他，这是你的责任 2、lua way 可以在 lua侧包装一下返回的类，再使用 ffi.gc 来注册析构函数

## 最佳实践
- 很多 lua 的内置函数都是全局变量，把它注册到本地来使用，性能会更好。
- 注意 Openresty 当中请求域名时会使用 Nginx 配置的 Dns 服务器，搜索 `resolver` 了解更多细节（Nginx 实现了一套内置的 DNS 解析）
- Openresty 中默认不读取 body ，可以通过以下方式打开
```
http {
    server {
        listen    80;

        # 默认读取 body
        lua_need_request_body on;

        location /test {
            content_by_lua_block {
                local data = ngx.req.get_body_data()
                ngx.say("hello ", data)
            }
        }
    }
}
```
或者局部开启
```
ngx.req.read_body()
```
- Openresty 有时 `ngx.req.get_body_data()` 读取不到数据时因为已经被转储到文件了，还需要从 `ngx.req.get_body_file` 中读取
- Nginx 有两个比较关键的参数 [client_body_buffer_size](http://nginx.org/en/docs/http/ngx_http_core_module.html#client_body_buffer_size) 和 [client_max_body_size](http://nginx.org/en/docs/http/ngx_http_core_module.html#client_max_body_size) ，前者控制缓冲区大小，大于这个部分的请求体会被转储为临时文件（文档上说明为部分或全部，但目前观测到的情况一般都是全部转储）；后者控制能接受的最大请求体，大于会被拒绝( 413 Request Entity Too Large)