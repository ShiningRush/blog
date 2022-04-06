+++
title = "golang context的正确打开姿势"
date = "2022-02-07T13:55:31+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["go"]
keywords = ["go","context"]
description = ""
showFullContent = false
readingTime = false
+++

# golang context的正确打开姿势
工作中发现很多同学不太了解 context 什么时候用，怎么用，这里从我的工作经验上简单分享下，golang中该如何优雅地处理错误。

## context目的
两个重要作用：
- 存活检测: 当前的工作上下文是否依然有效，工作上下文在不同的环境请求可能不同，比如在API服务中，上下文的生命周期等于请求的生命周期，在后台任务中，上下文的生命周期等于单次任务执行的生命周期。
- 存取值: 负责在上下文中传递公共参数，比如 调用链路、事务、操作用户的信息等

golang 对 Context 的定义是用于承载 `request-scoped values`, `cancellation signals` 和 `deadlines`，因此 `thread-local storage` 仅仅只是它的部分目的（由于 golang 屏蔽了线程因此这里称为 `goroutine-local` 会更合适些），

## goroutine-local 的优缺点

为什么没有 `goroutine-local`:
- gotoutine 比起线程更加轻量，通常我们会预期程序当中有大量的 goroutine 在活动，如果引入这个模式会导致性能内存都受到很大挑战，如何gc，如何传播，锁竞争激烈等
- 本身 thread-local 由于其灵活性也存在一些弊端，比如可能会破环你的层次结构，你可以在 dao 中获取到一个 http 的对象，这可能会带来意外的依赖
- 带来内存泄漏，比如在池化的场景中，如果 `local storage` 的内容没有及时擦除，那么它可能会遗留到下一个环境中

## 小结
总的来说：
- 如果你在开发一个新功能，并且它涉及到 `IO密集型` 的操作，比如读取数据库、发起网络等，目前还是建议将第一个函数设置为 context，它可以用于判断当前的上下文存活性、变量传递等等；如果不涉及，你可以传入 context，也可以直接传入具体的值
- 如果你在改造遗留项目，你可以改变它的参数，如果这个工作量非常巨大，那么可以考虑使用一些开源项目来实现 goroutine-local，比如 [gls](https://github.com/jtolio/gls), [routine](https://github.com/go-eden/routine)，但是这样做可能来带几个风险：
  + 性能损耗：几乎所有的开源库都存在这个问题，尤其是从 `stack` 当中读取的，还有 gc，锁竞争的代价 等
  + 意外的bug：内存泄漏，local storage 在各个goroutine中的传播等

