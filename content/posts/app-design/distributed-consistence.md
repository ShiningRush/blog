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
- Paxos: 共识算法，实现最为复杂，但是最为全面。
- Raft( ETCD, TiDB ): 基于 Paxos 的简化。
- ZAB( Zookeeper ): 基于 Paxos 的简化，类似 Raft。
- Gossip(Redis, consul, cassandra): 共识算法，利用节点传播来达到一致性。基于六度分隔理论（Six Degrees of Separation）哲学的体现，简单的来说，一个人通过6个中间人可以认识世界任何人。是一种通过传播消息来冗余容错的最终一致性算法。
- Bully: 选举算法

Paxos, Raft, ZAB 都很类似，通过选主之后通过Leader写入消息，再同步到Follower来保证强一致性，而Gossip通过传播来达到最终一致性。因此业界强调高可用的组件会使用 gossip 而强一致性会使用 其他三种。
还有单纯的主从模式，有时也会搭配这些算法一些使用。


## Paxos

