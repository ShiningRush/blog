+++
title = "分布式事务笔记"
date = "2021-02-12T15:05:42+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["分布式事务"]
keywords = ["分布式事务"]
description = ""
showFullContent = false
+++

# 分布式事务笔记

分布式事务的常见解决方案：
- 2PC: perpare + commit
- 3PC: can_commit + pre_commit + do_commit
- TCC: try + confirm + cancel
- Saga: Events/Choreography or Command/Orchestration
- 本地消息表: local + event

## 2PC

2PC 由协调器通知各个事务参与者，由 `perpare` 与 `commit` 两个阶段组成，前者通知所有参与者预备事务，后者通知参与者真实提交事务。如果有参与者在 `prepare` 阶段失败，那么会通知所有参与者回滚。

优点:
- 步骤少，简单

缺点：
- 协调器单点故障
- 会阻塞其他参与者的事务
- `commit` 阶段可能出现消息丢失

## 3PC
3PC 在 2PC 的基础上增加了各个参与者的超时时间，如果参与者在超时是时间内没有相应，那么视为失败。3PC 依然没有解决 `commit` 消息丢失的问题。

## TCC
TCC 相当于 3PC 的完善，除了超时机制外，还存在 `cancel` 步骤，意味着可以在 `commit` 失败后进行补偿。

## Saga
Saga 是一种长活事务的设计模式( Long-live transaction )，它基于 1987 年的一篇论文，主要有两种实现：Events/Choreography or Command/Orchestration
可以参考：[saga-pattern-implement-business-transactions-using-microservices-part](https://blog.couchbase.com/saga-pattern-implement-business-transactions-using-microservices-part/)、
[saga](https://docs.microsoft.com/en-us/azure/architecture/reference-architectures/saga/saga)。

这里说一些需要注意的点：
- `Choreography` 与 `Orchestration` 都有编排的意思，但是前者是不需要 `协调器(Orchestrator)` 的角色，都是参与者通过 `事件` 相互协作。参考 [choreography-vs-orchestration](https://medium.com/ingeniouslysimple/choreography-vs-orchestration-a6f21cfaccae)
- 微软的 [基于容器化的微服务架构设计](https://docs.microsoft.com/zh-cn/dotnet/architecture/microservices/architect-microservice-container-applications/asynchronous-message-based-communication) 极力推荐服务间通过异步消息来通信，其实就是 `Events/Choreography` 模式的实现。
- `Command/Orchestration` 模式可以将 `协调器` 角色集成到调用者服务中，也可以共用一个 `中心协调器`，一般更推荐前一种模式，除非业务中存在大量类似需求。两者都需要注意多个实例时可能会存在多次执行的问题，所以一般都会涉及选主操作。
- `Command/Orchestration` 模式下如果将 `协调器` 集中到调用者服务时，需要考虑多个副本间的竞争关系。如果使用 `k8s` 集群，可以通过 `Sidecar` 方式来进行 Leader选举 从而轻松解决竞争问题，参考 [contrib/election](https://github.com/kubernetes-retired/contrib/tree/master/election)。
- `Events/Choreography` 应该作为系统首选方案，只有事务涉及服务太多( > 4 )的情况再考虑 `Command/Orchestration`。
- `Events/Choreography` 要注意推送事件失败时的重试与补偿机制，需要的情况下可以使用 `本地消息表` 模式
- `Saga` 有两种错误恢复机制： `BackwordRecovery` 与 `ForwardRecovery`，前者指当某个参与者失败时，逆序对所有参与者调用补偿操作，这意味着这个事务的所有参与者都要实现一个补偿操作；而`ForwardRecovery` 则会重新对参与者发送请求，直到其成功，所以要求操作必须 `幂等`，由于这种恢复机制会假设永远成功，所以不需要实现补偿操作。
- 常用搭配：
  + `Events/Choreography` + `ForwardRecovery`: 实现最为简洁，应对事务复杂度较低，可以容忍实时性下限过低的场景，也可以避免由于事件补偿所导致的循环依赖问题
  + `Command/Orchestration` + `BackwordRecovery`: 实现最为复杂，应对复杂事务与无法容忍实时性下限过低的场景

## 本地消息表
这是 `eBay` 的一种设计模式，用于确保消息推送时的一致性。
通常情况调用服务很难保证 `提交自身事务` 与 `动作执行` 两者是强一致的，你很难保证网络不出问题，所以这种情况可以采用 `本地消息表` 模式，在本地存储中新建一个消息表，将需要发送的消息与事务数据放入同一个事务中落库，再异步发送消息，这样可以保证两个动作的强一致。这里值得一提是消息表可能不需要一张新的表，比如你可以在目标资源的表中新加一个字段来标识该资源的消息尚未被处理。

可以把这种模式作为 Saga 的一个补充，虽然我觉得多数情况有重试和补偿机制已经足够了。另外值得一提的是，在微软的架构指南中，提出了一个很有意思的设计，那就是处理第一个请求时不需要做任何业务操作，如果是 事件编排则发起消息，如果是命令编排则保存命令。

## 小结
阿里开源了他们的分布式事务框架 Seata( 阿里的 Fescar 和蚂蚁的 TCC方案做了整合 )，就从这个框架来看看上面几个理论的实践， seata 支持 AT(类似2PC), TCC, SAGA, XA(X/OPEN定义的规范)：
- AT: 需要支持 ACID 事务的存储，AT 模式通过解析SQL而生成 before 和 after 的镜像文件来实现补偿，用户不需要实现补偿动作
- TCC: 只需要实现 prepare, commit, rollback 三个动作即可，所有动作都需要用户自行实现
- SAGA： 不需要持有事务和锁，很适合长活事务，但是隔离性差，比如之前的结果被回滚了，在回滚之前是有可能被其他 client 获取到这个状态的。注意其实 seata 实现的是 saga 的命令编排模式，seata 作为一个中心协调器的角色
- XA：需要实现了XA协议的存储，流程和AT相同，为什么需要这个模式，是因为现存的三个模式都是基于补偿性事务实现的。
- 命令编排 + BackwardRecovery 进行补偿时需注意的地方：
  + 空回滚：由于超时而执行回滚的动作，实际 commit 没有发生
  + 悬挂：由于