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
// 这个非常有有意思，考虑正常 channel 的用法，也许第一直觉我们是使用上面的方式
// 但是多数开源项目用的都是下面这种方式，你是不是会好奇为啥 close 也会触发 select
// 其实这个可以理解为：如果select 了一个已关闭的 channel 会发生什么？
// 答案是：立即触发这个 case ，并且 ret, ok := closeCh 中，ok 为 false
// 这个特性配合 select-default 可以对 channel 的存活进行检测（但是不推荐做这样的事）
close(aChan)
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

使用信号的方式非常优雅，但是性能损耗较大，比后者慢了50%

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

## defer 特殊性
defer 先进后出性质应该是非常常见了，但是在返回值的地方有个特殊处理可能很少人注意到：
```go
func f1() (result int) {
    defer func() {
        result++
    }()
    return 0
}

func f2() (r int) {
     t := 5
     defer func() {
       t = t + 5
     }()
     return t
}
```

f1 的执行结果是 1, f2 的执行结果是 5，这里需要注意的是 `return` 语句并不是原子执行的，它是按照以下顺序来执行：
- 返回值赋值
- 插入 defer
- 函数返回

也就是说上面的 f1 例子中，是按照以下顺序：
- result = 0
- result ++
- return

f2同理。

## math/rand 与 crypto/rand 区别
两者区别在于前者是运算而得，后者从机器上获取，后者性能差异较大（数量级的差异）

## 后台 Goroutine 相关
我们经常会在一个程序中跑很多 goroutine 来执行一些周期性任务，这里会涉及到两个常见的动作：
- 同步：如果程序关闭时某些任务正在运行，我们应该尽量等待它运行完毕，此时可以使用 `sync.WaitGroup` 来做同步，然后组件暴露一个 `Close` 方法让用户显式调用。
- 泄露：如果你的组件是一个可能在程序的生命周期内反复创建又销毁的，那么要注意在组件不再使用后要停止这些 goroutine，你可以像前面一个场景一样暴力一个 `Close` 来解决，但是如果你不需要做同步，还有另一个简单而优雅的办法 `runtime.SetFinalizer` ，它可以在对象被GC时执行，你可以在这个函数中回收后台 goroutine，[go-cache]() 和 istio 的 lrucache 中都使用了这种方法来避免泄露。使用该函数有几点注意
  + 避免在终结器中添加循环引用，GC可以识别正常的循环引用，但是在终结器里面的无法处理
  + 这会造成对象生命周期拉长，因为终结器执行后会将对象的终结器清空，且留到下一个GC才会回收，大并发场景下不太友好
  + 由于终结器只会在对象被 GC 时回收，因此要注意在程序结束时不会触发的，所以不要用它来做一些持久化的工作，还需要注意如果有后台 goroutine 引用了对象也会导致它无法进入不可达状态，而无法触发终结器，这种时候可以使用一个 wrapper 对象将它包住，而将 wrapper  暴露给用户使用。


## sync.Cond 与 channel
查看 client-go 代码时发现他的限速队列用了 `sync.Cond` 看了下作用其实就是阻塞若干 goroutine 然后可以单个 or 整体唤醒(signal or broadcast)，有点好奇为啥有了 channel 还需要这个，然后就发现 go 有个[issue](https://github.com/golang/go/issues/21165) 在讨论是否要在 go2 中取消这个功能, 我总结了下：
- sync.Cond 的 signal 和 broadcase 都可以使用 channel 的 send 和 close 来替代
- sync.Cond 的性能更好（但是没看见有对比用例）
- sync.Cond 的语义符合 c 的 pthread 原语
- 只保持 channel 在用法上会更简洁
- sync.Cond 的用法较为繁琐，用不好容易出bug

## go mod 版本
go mod 正常会获取最新的 tag 作为默认版本，要求 tag 遵从语义版本 (vx.0.0-xxx)。
另外为了针对内部包快速迭代又不想打 tag 情况，提出了 伪版本的概念(pseudo-version)，格式如 vx.0.0(基础语义版本) - yyyymmddhhmmss(修订版本的UTC创建日期) - abcdefabcdef(VCS的commit hash 前缀)
伪版本不要手工输入可以使用
```
go get xxxxx@branch
go get xxxxx@commit hash
go get xxxxx@tag
``` 

来获取指定版本, 值得注意的是如果一个仓库使用了有效的语义版本作为 tag 且主版本不是 v0 v1, 那么在 go.mod 清单中它会在后面追加一个 "+incompatible" 的元数据标识

## go http lib
- http.Client http.RoundRrip