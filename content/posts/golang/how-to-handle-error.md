+++
title = "如何优雅地处理Error"
date = "2022-02-24T12:04:19+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["", ""]
keywords = ["", ""]
description = ""
showFullContent = false
readingTime = false
+++

# 如何处理Error
怎么处理Error是golang中一个非常关键的事情，因为golang的设计导致代码中到处都是类似以下的代码
```golang
if err != nil {
    ...
}
```

如果处理不得当，会导致代码膨胀得非常快且难以维护，比如：
```golang
if err != nil {
    metrics.Emit(...)
    log.Print(...)
    event.Emit(...)
    ...
    ...
    return err
}
```
面对这样的代码，可能错误处理所占的代码行数都会多于逻辑代码，显然不是我们愿意看到的。
上面描述的代码膨胀现象只是错误处理中常出现的问题之一，接下来我们聊聊日常开发中该如何优雅地处理错误。

## 错误的处理方式
错误的传递方式无非两种：
- 返回
- 不返回

选择返回or不返回的场景无非几种：
- 当函数位于顶层，比如 `API的接入层`、`conumer的handle` 等，此时无法返回
- 发生的错误是 `致命的`，会影响到整个程序的运行，此时应该抛出panic，阻止程序发生更加不可控的事情
- 发生的错误是 `预期的`，比如查重动作中查询数据库的数据不存在时，得到了一个 `NotFound` 的错误
- 其他情况均应该返回错误

而面对错误发生时，常见有以下的处理方式：
- 记录(日志、metrics、事件)
- 做一些业务逻辑
- 直接返回

这里只对 `记录` 单独展开说下，其他两项暂时没什么需要注意的。

### 记录
记录的常见手段有三种：`log`, `metrics`, `event`。
通常来说我们只需要记录其中一种即可，有的同学可能会有疑问，明明这三种不同的技术都有不同的适用场景，为什么说通常了一种就够了。
首先我们要聊聊这三个记录的核心目标是什么：
- `log`: 最传统的形式，是为了在程序运行时留下可追溯的信息，来辅助人类了解程序发生了什么，一般都会进行分级处理。
- `metrics`: 由于不同的 metrics 体系，实现不同，这里只概述下，metrics是用于观测某个属性的趋势
- `event`: 用于敏感操作的审计 or 广播变更

我在过往的工作中，见过某些同学，在某个很关键的场景出错后采用了以下的做法：
- 用 log 来记录详细的报错信息
- 同时记录一下 metrics 用于度量这个服务的错误率
- 发出一个 event，用于记录本次的操作审计

看起来似乎很合理，各个输出用于不同的目的，但细细琢磨一下，有两个疑问：
- log/event 难道不能用于度量错误率吗？
- metrics/event 难道不能记录错误详情吗？

第一个问题：可以，其实在没有 metrics 这种记录方式之前，早就已经存在针对 log 的关键信息来进行统计的用法，曾经我们在阿里云上就使用对于分类日志（请求日志、程序日志）的筛选来度量服务的稳定性
第二个问题：可以，我们可以在 metrics/event 的 tag/label 中携带上错误的原因和上下文。

所以说推荐在记录上只选择其中一种来记录错误，但是针对 `不同种类的错误` 可以有不同的记录方式，如何处理不同的错误，可以在框架中间件 or 顶层函数去决策。
无论你做度量or告警，都建议只针对一类记录来做，这样有几个好处：
- 维护的代码量减少
- 告警、度量来源统一，不需要怀疑是否某个数据记录有问题

> Tips
> 如果使用 event or metrics 来记录错误时，也建议在里面封装一层打印一下日志，格式相对可以简单些，便于调试开发时查看

## 错误处理的一些技巧
对于错误我有一些技巧分享

### 对错误添加更多的上下文
我们最常见到的方式如下：
```go
if err := readDb(); err != nil  {
    return err
}

if err := writeDb(); err != nil  {
    return err
}

if err := rpcCall(); err != nil  {
    return err
}
```

但是有没有想过，如果你在日志里面发现一条错误日志，上面的 message 只有 `err: connect timeout` 时，你能区分是从上面的哪个代码块返回的错误吗？
所以我们需要给错误添加更多的上下文，比如这样：
```go
if err := readDb(); err != nil  {
    return fmt.Errorf("read failed: %w", err)
}

if err := writeDb(); err != nil  {
    return fmt.Errorf("write failed: %w", err)
}

if err := rpcCall(); err != nil  {
    return fmt.Errorf("call xxx failed: %w", err)
}
```

这样很好，但是问题又来了，如果我们 `逐层` 为每一个函数的每一个返回点都添加上下文，这个成本也很高，那么我们什么时候才需要去添加上下文？
答案是：`错误发生点` or `边界`。
`错误发生点` 很容易理解，比如你的接口有输入限制，当输入参数不满足预期时，你会返回错误。
```go
if input.Param != expectedVal {
    return fmt.Errorf("params should not be [%s], expected: %s", input.Param, expectedVal)
}
```
上面的例子，发生错误时很明显能够知道错误发生在哪里，错误详情是什么，因此定位是很简单的。

`边界` 指项目代码与非项目代码的边界，比如你调用了非项目的SDK：gorm，web 等，这里包含标准与非标准的库，此时由于我们没法感知到 `错误发生点`，因此只能在边界上添加上下文，以便在错误发生时能够找到项目内的相关代码片段。
情况类似刚才举例的 conntect timeout。 

### 巧用 defer
在顶层函数中，我们必须要消化掉错误，此时我们很有可能会这样做
```go
func handle() {
    if err := methodA(); err != nil  {
        logs.Println(err)
        return
    }

    if err := methodB(); err != nil  {
        logs.Println(err)
        return
    }

    if err := methodC()); err != nil  {
        logs.Println(err)
        return
    }
}
```
此时我们可以用 `defer` 来统一处理
```go
func handle() {
    var err error
    defer func() {
        if err != nil {
            logs.Println(err)
        }
    }

    if err = methodA(); err != nil  {
        return
    }

    if err = methodB(); err != nil  {
        return
    }

    if err = methodC()); err != nil  {
        return
    }
}
```

### 错误分类
如果我们需要针对不同的错误场景做不同的处理，不要在各个函数写死，还是应该遵从在顶层函数处理的原则，但是我们可以返回项目自定义的错误，比如：
```go
type BaseError struct {
    Code int
    Msg string
    Context map[string]string
}

func (err *BaseError) Error() string {
    return fmt.Errorf("msg: %s, ctx: %+v", msg, err.Msg, err.Context)
}

var NotFoundError = &BaseError{Code: 404, Msg: "record not found"}
```

一来可以针对不同错误去定义不同的处理策略，二来我们可以更容易地区分出上下文，而不是直接耦合在 Msg 当中：
```go
func errHandle(err Error) {
    switch err.(type)
    case *BaseError:
        xxxx;
    default:
        xxxx;
}
```