+++
title = "分布式一致性笔记"
date = "2021-02-12T15:04:11+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["app-design","分布式一致性"]
keywords = ["分布式一致性"]
description = ""
showFullContent = false
+++

# 分布式一致性

目前主流算法如下：
- Paxos: 共识算法，实现最为复杂，但是最为全面。又分为 Basic-Paxos 与 Multi-Pax，用于处理一条提案与一批提案的区分。
- Raft( ETCD, TiDB ): 基于 Paxos 的简化，是在几个限制前提下的 Paxos: 发送的请求的是连续的, 也就是说Raft的append 操作必须是连续的， 而Paxos可以并发 (这里并发只是append log的并发, 应用到状态机还是有序的)。Raft选主有限制,必须包含最新、最全日志的节点才能被选为leader. 而Multi-Paxos没有这个限制，日志不完备的节点也能成为leader。Multi-Paxos允许并发的写log,当leader节点故障后，剩余节点有可能都有日志空洞。所以选出新leader后, 需要将新leader里没有的log补全,在依次应用到状态机里。参考 [TIDB 架构及分布式协议Paxos和Raft对比](https://blog.51cto.com/liuminkun/2377029)
- ZAB( Zookeeper ): 基于 Paxos 的简化，类似 Raft。
- Gossip(Redis, consul, cassandra): 共识算法，利用节点传播来达到一致性。基于六度分隔理论（Six Degrees of Separation）哲学的体现，简单的来说，一个人通过6个中间人可以认识世界任何人。是一种通过传播消息来冗余容错的最终一致性算法。因为会耗费大量带宽，因此使用该算法的组件基本都是用它来传播元数据(比如redisc cluster, consul都用其来传播节点信息)，以实现去中心化的管理。同时它也有三种实现模式：DirectMall, Anti-Entropy, Rumor-Mongering，DirectMall优势在于尽可能多的通知周围节点，同时发送全量信息，后两者区别在于每次传播到底是全量的还是仅传播增量的数据。参考文档：[Gossip Protocols](http://www.cs.cornell.edu/courses/cs6410/2016fa/slides/19-p2p-gossip.pdf), [Gossip 协议](https://zhuanlan.zhihu.com/p/41228196),[https://zthinker.com/archives/%E6%BC%AB%E8%B0%88gossip%E5%8D%8F%E8%AE%AE%E4%B8%8E%E5%85%B6%E5%9C%A8rediscluster%E4%B8%AD%E7%9A%84%E5%AE%9E%E7%8E%B0](https://zthinker.com/archives/%E6%BC%AB%E8%B0%88gossip%E5%8D%8F%E8%AE%AE%E4%B8%8E%E5%85%B6%E5%9C%A8rediscluster%E4%B8%AD%E7%9A%84%E5%AE%9E%E7%8E%B0)
- Bully: 选举算法

Paxos, Raft, ZAB 都很类似，通过选主之后通过Leader写入消息，再同步到Follower来保证强一致性，而Gossip通过传播来达到最终一致性。因此业界强调高可用的组件会使用 gossip 而强一致性会使用 其他三种。
还有单纯的主从模式，有时也会搭配这些算法一些使用。


## Paxos

