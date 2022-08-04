+++
title = "Nginx 配置最佳实践"
date = "2022-06-21T19:26:31+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["", ""]
keywords = ["", ""]
description = ""
showFullContent = false
readingTime = false
+++

# Nginx 配置最佳实践

## Buffer
注意Nginx作为代理服务器时的几个关键行为：
1. 默认会缓存 request 的 body，直到接收完所有body后才会转发请求到后端
2. 默认会缓存上游服务的 response，直到接收完所有 response 的body or body超过了设定值才会将请求转发给后端


这两个行为的目的都是为了克服慢client带来的影响，比如 client 在发起请求时很慢，大量的连接会贯穿到 upstream，而接收响应也是类似原理，nginx期望尽可能快地接收完 upstream 响应，以释放它的相关资源。
现在我们来详细看看过程中涉及哪些参数，以一个请求发起为例

1. client 发起请求，nginx 根据 `proxy_request_buffering` 判断是否需要缓存请求 body，默认为 `on`，即缓存请求
2. 如果开启缓存，那么 nginx 会不断缓存 body 到内存中，大小为 `client_body_buffer_size` 所指定的大小，默认为一个内存页大小 4k(32bit os)/8k(64bit)，不过这个参数不一定适合现在的硬件了，参考 [这里](https://draveness.me/whys-the-design-linux-default-page/)，如果超过这个大小，则会写入临时文件，路径为 `client_body_temp_path` 所配置。
3. 1,2步骤完毕后，开始向后端转发，然后等待后端请求返回
4. 开始接收 response ，根据 `proxy_buffering` 判断是否需要缓存响应 body，默认为 `on`，即缓存响应
5. 如果开启缓存，和请求类似的是也有一个 `proxy_buffers` 来决定缓冲区大小，默认 `8 4k`(8 * 4k大小)，当body大于这个buffer时会转储本地文件，文件大小为 `proxy_max_temp_file_size` 限制，每次写入大小为 `proxy_temp_file_write_size`，默认是 `proxy_buffer_size` 和 `proxy_buffers` 之和，当大于临时文件限制时将转为同步传输，同步传输将使用 `proxy_buffers` 定义的空间作为 buffer，同时还有一个 `proxy_busy_buffers_size` 用于控制 buffer 中的哪一部分用于回传给 client,因为传输的过程中，整个 buffer 都将被标记为 busy，不可用，因此为了实现边接边回传的均衡，建议 `proxy_busy_buffers_size` 不要大于 `proxy_buffers` 的一半。

如果关闭 nginx 的 buffering，那么nginx将使用 `proxy_buffer_size` 配置的 size 作为 buffer 来传输文件，可以通过提升这个值来加快传输速度

了解以上机制后，我们可以得出以下最佳实践：
1. 上传文件时可以关闭 buffer ，避免请求写入磁盘导致磁盘IO，也更节约内存
2. 下载文件时可以继续保持 buffer，但是调整 `proxy_max_temp_file_size` 为0，避免磁盘IO带来的性能劣化，同时要保持 `proxy_busy_buffers_size` = `proxy_buffers` / 2，避免退化为串行传输。（遗憾的是，上传时没有这个配置。）
3. 通常API访问时要注意调整 buffer，避免磁盘IO

参考：
- [Avoiding the Top 10 NGINX Configuration Mistakes](https://www.nginx.com/blog/avoiding-top-10-nginx-configuration-mistakes/#configuring-buffers)
- [Performance Tuning – Tips & Tricks](https://www.nginx.com/blog/performance-tuning-tips-tricks/)
- [Tuning proxy_buffer_size in NGINX](https://www.getpagespeed.com/server-setup/nginx/tuning-proxy_buffer_size-in-nginx)
- [Nginx反向代理缓冲区优化](https://www.cnblogs.com/xiewenming/p/8023090.html)