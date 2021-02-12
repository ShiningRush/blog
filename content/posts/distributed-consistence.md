+++
title = "分布式一致性笔记"
date = "2021-02-12T15:04:11+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["分布式一致性"]
keywords = ["分布式一致性"]
description = ""
showFullContent = false
+++

# 分布式一致性

目前主流算法如下：
- Paxos: 共识算法，实现最为复杂，但是最为全面
- Raft( ETCD, TiDB ): 基于 Paxos 的简化
- ZAB( Zookeeper ): 基于 Paxos 的简化，类似 Raft
- Gossip: 共识算法，利用节点传播来达到一致性
- Bully: 选举算法


## Paxos

