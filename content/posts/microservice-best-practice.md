+++
title = "微服务最佳实践"
date = "2021-02-12T15:03:15+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["微服务"]
keywords = ["微服务"]
description = ""
showFullContent = false
+++

# 微服务最佳实践
本文记录了工作当中积累的一些经验。

## 一定要做到的点
- 每个请求一定要有 `request-id`，一般由调用方生成，如果你是在推动一个遗留项目，那么可以在网关或者`AOP`层去生成它，这样才能够在繁杂的日志中找到特定的请求信息。建议由 `uuid` 或者 雪花算法生成。
- 请求日志和错误日志分离，这么是因为这两者具有完全不同的关注点，且前者格式相对固定。分离他们可以在清洗日志时使用不同的格式。
- 如果你们的系统引入了调用链追踪，那么可以使用 `request-id` 作为 `trace-id`。
- 对于每一个资源的操作都应有反馈，比如删除不存在的资源或者创建同名的资源，我们都应该返回特定的错误码，通常可以使用三类: `NotFoundError`, `ConflictError`, `InternalError`

## 框架选择
- 选择什么框架其实不是很重要，建议把代码构建成框架无关的形式，一般情况下可以使用一些抽象中间件完成

## 分布式原语
可以参考 Bilgin Ibryam 的博客 [Top 10 must-know Kubernetes design patterns](http://www.ofbizian.com/2020/05/top-10-must-know-kubernetes-design.html)