+++
title = "api压力测试工具简介"
date = "2021-02-12T14:48:31+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["benchmark", "wrk", "vegeta"]
keywords = ["benchmark", "wrk", "vegeta"]
description = ""
showFullContent = false
+++

# api压力测试工具简介
这里介绍两个我常用的测试工具：vegeta 和 wrk，现在我一般都用前者

## Vegeta( 贝吉塔 )

这是一个开源软件，类似 wrk ，但是是用 go 写的，功能更为强大，也更易用一些。

[Github 地址](https://github.com/tsenart/vegeta)

下面是命令列表：
```
Usage: vegeta [global flags] <command> [command flags]

global flags:
  -cpus int
    	Number of CPUs to use (defaults to the number of CPUs you have)
  -profile string
    	Enable profiling of [cpu, heap]
  -version
    	Print version and exit

attack command:
  -body string
    	Requests body file
  -cert string
    	TLS client PEM encoded certificate file
  -chunked
    	Send body with chunked transfer encoding
  -connections int
    	Max open idle connections per target host (default 10000)
  -duration duration
    	Duration of the test [0 = forever]
  -format string
    	Targets format [http, json] (default "http")
  -h2c
    	Send HTTP/2 requests without TLS encryption
  -header value
    	Request header
  -http2
    	Send HTTP/2 requests when supported by the server (default true)
  -insecure
    	Ignore invalid server TLS certificates
  -keepalive
    	Use persistent connections (default true)
  -key string
    	TLS client PEM encoded private key file
  -laddr value
    	Local IP address (default 0.0.0.0)
  -lazy
    	Read targets lazily
  -max-body value
    	Maximum number of bytes to capture from response bodies. [-1 = no limit] (default -1)
  -max-workers uint
    	Maximum number of workers (default 18446744073709551615)
  -name string
    	Attack name
  -output string
    	Output file (default "stdout")
  -proxy-header value
    	Proxy CONNECT header
  -rate value
    	Number of requests per time unit [0 = infinity] (default 50/1s)
  -redirects int
    	Number of redirects to follow. -1 will not follow but marks as success (default 10)
  -resolvers value
    	List of addresses (ip:port) to use for DNS resolution. Disables use of local system DNS. (comma separated list)
  -root-certs value
    	TLS root certificate files (comma separated list)
  -targets string
    	Targets file (default "stdin")
  -timeout duration
    	Requests timeout (default 30s)
  -unix-socket string
    	Connect over a unix socket. This overrides the host address in target URLs
  -workers uint
    	Initial number of workers (default 10)

encode command:
  -output string
    	Output file (default "stdout")
  -to string
    	Output encoding [csv, gob, json] (default "json")

plot command:
  -output string
    	Output file (default "stdout")
  -threshold int
    	Threshold of data points above which series are downsampled. (default 4000)
  -title string
    	Title and header of the resulting HTML page (default "Vegeta Plot")

report command:
  -buckets string
    	Histogram buckets, e.g.: "[0,1ms,10ms]"
  -every duration
    	Report interval
  -output string
    	Output file (default "stdout")
  -type string
    	Report type to generate [text, json, hist[buckets], hdrplot] (default "text")

examples:
  echo "GET http://localhost/" | vegeta attack -duration=5s | tee results.bin | vegeta report
  vegeta report -type=json results.bin > metrics.json
  cat results.bin | vegeta plot > plot.html
  cat results.bin | vegeta report -type="hist[0,100ms,200ms,300ms]"
  vegeta attack -duration=5s -targets=./targets.http | vegeta report
```

 这里主要提一下 `-format=json` 时，可以使用类似以下的 json 文件作为请求源，这样动态能力会更强一些，在压测 TPS 时特别有用：
 ```json
{ "method":"GET", "url": "xxxx", "body":"base64(xxxxx)", "header": { "key": "value" } }
{ "method":"POST", "url": "xxxx", "body":"base64(xxxxx)", "header": { "key": "value" } }
 ```

如果没有特别的要求，一般压测 QPS使用 `-format=http`，这也是默认值。
```text
GET http://user:password@goku:9090/path/to
X-Account-ID: 8675309

DELETE http://goku:9090/path/to/remove
Confirmation-Token: 90215
Authorization: Token DEADBEEF
```


## Wrk 用法
### Wrk是什么
wrk是一个开源的http的压测工具，它封装了很多开源项目的，比如`redis`的`ae`( *一个事件循环的非阻塞网络库，底层封装了 epoll 和 kqueue*) 和 `nginx` 的 `http-parser`。基于这些优秀的开源项目，所以它的性能相当高。不过也因为它用了`epoll`和`kqueue`的原因，所以目前只支持`linux`平台。
同时还集成了`LuaJIT`，所以可以自己写Lua脚本，放在 `/scripts` 目录下。[点击这里](https://github.com/wg/wrk) 访问它的 Github。

### 安装
在Github的Release中下载最新版解压即可。
```
wget https://github.com/wg/wrk/archive/4.1.0.tar.gz ./
tar -xzvf ./4.1.0.tar.gz
make
```
安装后自行选择是否是否将二进制文件放入`$PATH`中

### 基础用法
```bash
wrk -t12 -c400 -d30s http://127.0.0.1:8080/index.html
```

这个命令代表起 `12` 个线程来保持 `400`个 http连接，持续`30`秒。
输出如下：
```bash
Running 30s test @ http://127.0.0.1:8080/index.html
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   635.91us    0.89ms  12.92ms   93.69%
    Req/Sec    56.20k     8.07k   62.00k    86.54%
  22464657 requests in 30.00s, 17.76GB read
Requests/sec: 748868.53
Transfer/sec:    606.33MB
```

### 参数解释
```bash
-c, --connections: 要保持的连接总数，每个线程要处理链接数 = 连接总数/线程数

-d, --duration:    测试持续事件, 如 2s, 2m, 2h

-t, --threads:     总线程数

-s, --script:      lua脚本, 参考上面的链接

-H, --header:      要追加的http header e.g. "User-Agent: wrk"

    --latency:     统计详细的延迟

    --timeout:     如果一个请求在该时间内没有返回，则记录一个超时
```

### 使用技巧
运行wrk的计算机必须具有足够数量的临时端口(Port)，并且关闭的端口应该快速回收。为了处理初始连接突发，服务器(listen(2))[http://man7.org/linux/man-pages/man2/listen.2.html]的backlog应该大于正在测试的并发连接的数量。
(*这里稍微解释一下, listen 指linux的系统函数, 它的第二个参数 backlog 指定了能够服务的客户端最大数量，如果超过这个数量的请求都会被拒绝掉。 *)
仅更改HTTP Mehod，Path，Header 或 Body 的不会对性能产生影响。每个请求的操作（特别是构建新的HTTP请求）以及使用response（）必然会减少可以生成的负载量。