+++
title = "数据库知识小记"
date = "2021-12-25T12:52:27+08:00"
author = "shiningrush"
authorTwitter = "" #do not include @
cover = ""
tags = ["db"]
keywords = ["db", "b+", "lsm", "hash"]
description = ""
showFullContent = false
+++

# 数据库知识小记
最近任务涉及到要在redis里面存储结构化数据，并且还要支持范围查询，同时方案设计时还需要考察各类数据库，因此调查了一下目前业界的数据库。
目前业界数据库种类繁多，根据使用方法不同可以分为：
- SQL: Mysql, PostgresQL, SQLServer, Oracle
- NoSQL: 不同传统的SQL数据库，又分成了：
  + 文档型数据库：CouchDB, MongoDB, ElasticSearch等
  + 图数据库：JanusGraph等
  + 键值数据库： ETCD, Redis, memcached, Zookeeper, Consul
  + 列族数据库：Cassandra, BigTable, HBase
  + 时序型数据库：InfluxDB, OpenTSDB, Prometheus, m3db等等
- NewSQL: Tidb


数据库太多，这里不一一枚举了，有比如kv下面有非常多的数据库，它们背后落盘的方式都不太相同，但是大致可以分为几种：
- HASH：键值类数据常用，优点是易维护，查询快，缺点是范围查询太慢
- B树：很多存储都用了B树的方式去构建数据，典型的比如Mysql（B+），Mongo(B-)
- LSM树：在机器磁盘上表现非常出色的写性能，用于降低存储成本，但是读上面牺牲比较大

同时在索引上也分为两类：
- 聚集索引(主索引)：索引上直接存储了数据，查到后直接返回数据
- 非聚集索引(辅助索引)：索引上仅保存了主键索引号，查到后仍然需要到主键索引上获取数据

## 参考文档
- [Wiki NoSQL](https://zh.wikipedia.org/wiki/NoSQL)
- [一篇文章讲透MySQL为什么要用B+树实现索引](https://cloud.tencent.com/developer/article/1543335)
- [cap](https://cloud.tencent.com/developer/article/1437772): 实际应用中的可用性和CAP可用性并不相同。你应用的可用性多数是通过SLA来衡量的（比如99.9%的正确的请求一定要在一秒钟之内返回成功），但是一个系统无论是否满足CAP可用性其实都可以满足这样的SLA。