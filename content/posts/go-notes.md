+++
title = "Go语言笔记"
date = "2021-02-12T14:40:39+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["go"]
keywords = ["go"]
description = ""
showFullContent = false
+++

# go语言笔记
这里记录一些与 go 相关的点

## GODBUEG
可以通过环境`GODEBUG` 来输出 go 程序的相关信息，可选命令如下：
- `allocfreetrace`
- `clobberfree`
- `cgocheck`
- `efence`
- `gccheckmark`
- `gcpacertrace`
- `gcshrinkstackoff`
- `gctrace`
- `madvdontneed`
- `memprofilerate`
- `invalidptr`
- `sbrk`
- `scavenge`
- `scheddetail`
- `schedtrace`
- `tracebackancestors`
- `asyncpreemptoff`

完整命令参考 [这里](https://golang.org/pkg/runtime/)
值得一提得是，可以通过`name=val,name=val`来启用多个命令，如：
- ``

## http库的 ServeHttp 不能修改 request 的内容
很有意思的一点，在 `ServeHttp` 代码注释中写着
> Except for reading the body, handlers should not modify the
> provided Request.

稍微调查了下原因，是因为大量的现存代码会受到影响，所以在 `http.StripPrefix` 函数中，对 `Request` 进行了深拷贝。
```go
func StripPrefix(prefix string, h Handler) Handler {
	if prefix == "" {
		return h
	}
	return HandlerFunc(func(w ResponseWriter, r *Request) {
		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			r2 := new(Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			h.ServeHTTP(w, r2)
		} else {
			NotFound(w, r)
		}
	})
}
```

参考资料: [net/http: allow handlers to modify the http.Request](https://github.com/golang/go/issues/27277), [stack-overflow](https://stackoverflow.com/questions/13255907/in-go-http-handlers-why-is-the-responsewriter-a-value-but-the-request-a-pointer?rq=1)

## 编译的二进制无法在 alpine 中运行
当项目引用了 `net` 包之后，因为网络库在不同平台下的实现不同，所以它默认依赖了 `cgo`来做动态链接，这里有几个解决办法：
- 禁用 `cgo`，如果你的项目没有依赖的话。`CGO_ENABLED=0 go build -a`，`-a` 表示让所有依赖库进行重编（这里验证过不带 -a 也可以正常工作，有点奇怪），禁用 cgo 后将会以静态连接的方式编译二进制包。
- 强制 go 使用一个特定的网络实现。`go build -tags netgo -a`
- 添加动态连接库，如果你的项目要依赖`cgo`
```
RUN apk add --no-cache \
        libc6-compat
```

这里简单提一下，还可以使用 `-ldflags="-s -w"` 来减少生成的二进制体积，它裁剪了程序的调试信息。
还有个`-installsuffix cgo`，go 1.10 以后已经不再需要。

参考 [这里](https://stackoverflow.com/questions/36279253/go-compiled-binary-wont-run-in-an-alpine-docker-container-on-ubuntu-host/36308464#36308464)

## goroutine 的性能小记
由于 goroutine 实现中包含了一把互斥锁，因此在多个 worker 竞争一个 channel 时会比较消耗性能，此时可以考虑对 channel 进行分片，实验发现以下两个配置效果比较好：
- 一个 worker 对应一个 channel，不过这样比较占用内存，同时也要注意数量不能无限递增，性能最高，但是空间占用最大
- 取CPU逻辑线程数(runtime.GOMAXPROCES) 作为分片数，这是在 m3 中看到的做法，效果不错，性能次之不过很节省空间

其他配置当然也可以，不过几乎和上面两种差不太多，前者比后者快了20%左右。做分片要注意考虑哈希算法的实现和性能影响，可选方案：
- atomic 自增 + 取余是一个不错的实现，但是注意自增的上限值。
- 如果有特征值，murmur3 + 取余是个完美的解决方案。

在优雅退出时，有几种方式实现：
- 利用 for + select + 退出信号
```go
// task
for {
	select {
	case <- aChan:
		doTask()
	case <- closeChan:
		end()
	}
}

// close
closeChan <- struct{}{}
```
- 利用 atomic + for
```go
// task
for isClose > 0 {
	<- aChan
	doTask()
}

// close
atomic.AddInt32(&isClose, 1)
aChan <- nil
```

使用信号的方式非常优雅，但是性能损耗比后者慢了50%

## 静态库
当你对一个 `package` 不是 main 的 go文件进行编译时，会得到一个 pkg 的归档文件 `x.a`，它可以提供给其他go引用，但是看不见声明与源码 

## 失败重试
常见的失败重试机制：
- 固定时间 (FixedDelay)
- 指数退避算法( Exponential Backoff )

这里简单介绍指数退避算法，即每次失败重试后，延迟的时间倍增，等待的时间为 2^n + `random_number_milliseconds`，比如：
失败第一次后，等待 1 + `random_number_milliseconds`
失败第二次后，等待 2 + `random_number_milliseconds`
失败第三次后，等待 3 + `random_number_milliseconds`

这里之所有会有一个随机的毫秒数，是为了防止通信双方如果在失败都尝试按照指数倍增的方式去重传，那么有可能会一直冲突，所以补充了一个随机值。

go 有个代码结构非常不错的重试库：[retry-go](https://github.com/avast/retry-go)

## 避免混用 GOPATH 和 go module
go 在编译时会优先使用 `GOPATH` 下面的源码，即便你的工作目录不处于 `GOPATH` 中，踩过一个坑就是一直在项目目录下编译，但是得到二进制文件都是没有改动过的，就是因为 `GOPATH` 下存在同名项目