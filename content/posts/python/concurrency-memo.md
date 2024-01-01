+++
title = "Python-并发小记"
date = "2023-04-23T14:48:43+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["python"]
keywords = ["", ""]
description = ""
showFullContent = false
readingTime = false
+++

# Python-并发小记
这里记录一些最佳实践，避免踩坑

## 如果使用 async/await，尽量保持项目统一
由于 thread 与 coroutine 的执行方式差异较大，在项目中混用时可能会出现问题，比如 coroutine 没有得到执行，以及 coroutine 的线程安全性
如果由于历史原因导致必须新技术混用，请确保好以下几点：
- `同步阻塞函数`可以使用 `eventloop.run_in_executor` 来包装为异步
- `同步阻塞函数`如果想要在其他事件循环（或者说是线程，因为默认是执行在运行线程上）上运行一个新的`异步函数`，使用 `asyncio.run_coroutine_threadsafe`
- 在进入 async/await 语法中后，请确保后续语法不要再混入同步函数
- 使用 `asyncio.run` 而避免直接去操作事件循环
- 在同步调用中使用 `asyncio`时尽量新起一个线程来run，不然可能会在上游出现已经运行中的事件循环从而导致冲突

## 如果项目有coroutine的代码，避免使用 threading.local
这是因为 threading.local 并不是协程安全的，请使用 `contextvars.ContextVar('var', default="default")` 来替代，否则可能出现变量的竞写导致预期外的情况

## 生成新的 thread 要注意 thread.local 与 context var 的传播
由于python不像是.Net 那样，默认会对子线程进行传播，因此在创建新线程or线程池时要注意 `thread.local` 和 `context var` 的显式拷贝（进程也一样）

## 尽量使用 asyncio 封装函数，不要直接接触evetloop
这是因为底层的机制比较复杂，比如对于 `main thread` 会默认初始化一个事件循环，但是新创建的子线程却不会。

## 注意避免使用同步IO
比如老牌的 `requests` 库，底层是基于 `urllib3` 的，很遗憾，这是个同步的 io 库，它会导致线程的阻塞而极大影响性能。
可以考虑 `aiohttp` or `httpx` 都是比较热门的异步 io 库，`httpx` 同时支持同步和同步的client语法，而 `aiohttp` 仅支持异步。
另外如果项目已经使用了 `request` 需要注意以下几点：
- Session 如果不复用，那么你将无法使用连接池
- 连接池默认最大连接为1，如果指定了`pool_block` 为 true，那么会导致client阻塞，这个可以参考 urllib3 的 [配置页面](https://urllib3.readthedocs.io/en/stable/reference/urllib3.connectionpool.html)

## 小心eventloop
python只会为主线程自动生成一个loop，因此如果你在子线程中调用 `get_event_loop` 会得到错误 **RuntimeError: There is no current event loop in thread**。
参考 [runtimeerror-there-is-no-current-event-loop-in-thread-in-async-apscheduler](https://stackoverflow.com/questions/46727787/runtimeerror-there-is-no-current-event-loop-in-thread-in-async-apscheduler)
正确的做法是为这个线程重新初始化一个：
```python
loop = asyncio.new_event_loop()
asyncio.set_event_loop(loop)
```