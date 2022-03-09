+++
title = "DDD实践心得"
date = "2021-02-12T14:37:29+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["app-design","领域驱动开发"]
keywords = ["DDD","领域驱动开发"]
description = ""
showFullContent = false
+++
# DDD实践心得

## ApplicationService vs DomainService
- 两者的相同处：都是对领域对象的组装。
- 两者不同处：`ApplicationService`不应该包含业务逻辑。
- 如何判断是否包含业务逻辑：当一个`Service`中包含了分支(复杂度大于1), 且分支条件由函数本身决定，那么这个`Service`包含了业务逻辑
- DomainService是否纯粹: 当一个`DomainService`仅仅包含领域对象时，则为纯粹的领域服务
- 实体是否可以依赖`DomainService`: 非纯粹的领域服务最好不要注入，纯粹的领域服务谨慎注入, 服务多数情况下应该处于上层
参考[Domain services vs Application services](https://enterprisecraftsmanship.com/posts/domain-vs-application-services/)

## 实体是否可以注入仓储
不可以，有以下几点原因
- 领域对象的知识应该尽可能纯粹, 仅包含业务知识
- 实体不应该知道另一个实体如何保留的知识，这样会过多让实体的生命周期边界变得模糊
- 如果你的实体需要访问仓储, 可能这段逻辑更适合放入`DomainService`中


## 杂项
- 仓储只应该提供 `增删改` 相关逻辑的查询，不要直接混用提供给调用方和提供给内部逻辑的 `List` 查询方法，因为对外一般都需要分页等等，但内部不需要
- 怎么识别聚合根？考虑聚合根的目的在于区分边界，边界内的实体都应该共享一个生命周期。那么最简单的一个法则就是：某批实体必须依托于另一个实体才有价值，那么另一个实体可以和它组成聚合。
- 对于复杂查询，首先应该考虑统一语言的建模中能否自然地包含它，如果不行，那么可以参考CQRS模式，建立一个专门做query的service，不过你应该时刻关注这个servic如何被消费了，以免产生误用
