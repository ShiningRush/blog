+++
title = "Leader选举机制"
date = "2021-02-12T14:45:18+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["选主","k8s"]
keywords = ["选主","k8s"]
description = ""
showFullContent = false
+++

# Leader选举机制
## 简介
Leader选举机制是为了解决分布式系统中一致性问题常见的一种机制，目前多数开源组件的选举算法都来自`Paxos`或者其简化版本`Raft`。

## 算法简介
- [Paxos](https://zh.wikipedia.org/wiki/Paxos%E7%AE%97%E6%B3%95)
- [Raft](https://raft.github.io/): 是对 `Mul-ti Paxos` 的一个改造，论文将算法分为三个部分：`Leader election`, `Log replication`, Safety.选举只是其中一部分，在选举中，节点分为三个状态: `follower`, `leader`，`candidate` ，开始所有节点均为 `follower` 状态，在指定的 `election timeout(通常 150~300ms)` 到期后，为自己投票(VoteCount+1)，同时向其他节点发送 `RequestVote` 消息，其中包含当前的任期编号(Term)和自己的标识，为了避免节点同时发起 Vote 所以 election timeout 都是在`允许范围内随机化`，当节点获得多数的票数，则立即变为 Leader，同时开始定期(heartbeat timeout)发送 Append Entries 消息，其他节点受到则会刷新 election timeout，同时回复消息。其他节点在接受到 RequestVote 会判断当前任期是否已经处理，如果已经处理则不会回复任何消息(此时依旧会触发 election timeout)
- [Bully](https://en.wikipedia.org/wiki/Bully_algorithm): 使用一个ID来竞选 Leader，最大者胜出。
除了这些算法外，如果你有用到额外存储，其实可以利用存储来实现选主，原理和分布式锁类似，多个实例争抢一个 key，谁抢到谁是leader，同时约定好ttl，定期继任和竞争。

## k8s 中的实现
了解k8s的读者应该知道，k8s早期的`master`是没有实现高可用的，所以早期社区出现各种不同的考可用方案，直到v1.13后官方才原生支持高可用。
它的实现方案就是使用了 `Leader选举`，并且利用了`k8s`中`endpoint`来优雅地完成云原生的选举机制。
参考[client-go leader-election](https://github.com/kubernetes/client-go/blob/b8fba595e8fa8e1f8dbad9b31129da74b3b6466b/tools/leaderelection/leaderelection.go#L76) 

另外还有一个SideCar使用这个库来完成 [contrib/election](https://github.com/kubernetes-retired/contrib/tree/master/election).