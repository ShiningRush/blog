+++
title = "数据库知识小记"
date = "2021-12-25T12:52:27+08:00"
author = "shiningrush"
authorTwitter = "" #do not include @
cover = ""
tags = ["db"]
keywords = ["db","b+","lsm","hash"]
description = ""
showFullContent = false
+++

# 数据库知识小记
最近任务涉及到要在redis里面存储结构化数据，并且还要支持范围查询，同时方案设计时还需要考察各类数据库，因此调查了一下目前业界的数据库。
目前业界数据库种类繁多，根据使用方法不同可以分为：
- SQL: Mysql, PostgresQL, SQLServer, Oracle
- NoSQL: 不同传统的SQL数据库，又分成了：
  + 文档型数据库：CouchDB, MongoDB, ElasticSearch等
  + 图数据库：JanusGraph, HugeGraph等
  + 键值数据库： ETCD, Redis, memcached, Zookeeper, Consul
  + 列族数据库：Cassandra, BigTable, HBase
  + 时序型数据库：InfluxDB, OpenTSDB, Prometheus, m3db等等
- NewSQL: Tidb


数据库太多，这里不一一枚举了，有比如kv下面有非常多的数据库，它们背后落盘的方式都不太相同，但是大致可以分为几种：
- HASH：键值类数据常用，优点是易维护，查询快，缺点是范围查询太慢，节点的增删成本高，涉及到 Reshard 。
- B树：很多存储都用了B树的方式去构建数据，典型的比如Mysql（B+），Mongo(B-)，优点是增删成本低，缺点是查询速度不如 HASH, O1 vs Ologn
- LSM树：在机器磁盘上表现非常出色的写性能，用于降低存储成本，但是读上面牺牲比较大，和B树类似

同时在索引上也分为两类：
- 聚集索引(主索引)：索引上直接存储了数据，查到后直接返回数据
- 非聚集索引(辅助索引)：索引上仅保存了主键索引号，查到后仍然需要到主键索引上获取数据，这个过程被称为 `回表`，如果 `查询的所有字段` 都属于辅助索引则不再需要回表，这个现象被称为 `覆盖索引(cover index)`，其原因是你用到的字段值都用于构建索引了，所以当然不需要回表查询其他数据

## 深度分页
所有存储都面临这个问题：`如果一大批数据进行非常靠后的查询，比如千万级别的结果集，跳到最后几条`，基本上现有存储的 skip 都是基于范围扫描实现的，比如 skip(100w).limit(1)，扫描 100w01条数据后抛弃前 100w 行记录来实现，因此 skip 太大后会导致大量的无效扫描。

这个问题有几个解法
- 调整业务：绝大多数这种场景都是不合理的，当你的默认条件能够查询出上百万条数据时，更合理的操作是缩小数据范围，而不是真的让用户翻到最后一页，这里从交互上有几个处理方法
  + 当超过预定的数据时(比如 1w)，只返回 1w，并提示用户缩小范围
  + 默认的搜索条件添加一些区分粒度大的条件，比如日期，可以有效控制默认页面的数据窗口
- 游标：游标就是在上次翻译的结果上做好标记，比如选取一个递增的字段，记录返回时的极值（通常选择主键ID），下次查询时携带上这个字段，然后用其作为查询条件来进行索引，而不是直接 skip。方案优点是易于实现，缺点是不支持随机翻页，只能查看下一页 or 上一页。
- 更换存储方案：通常这么大的统计需求一般都是建立 `数仓` 来实现了，不要直接基于 `OLTP` 的存储去实现。

## 隔离级别
先说说可能出现的几个问题：
- `脏读(dirty reads)`: A事务过程中读取到了B事务尚未提交的数据，这个是应该严格避免的。
- `不可重复读(non-repeatable reads)`: A事务过程中对相同行数据的两次查询返回结果不同，这是由于过程中有其他事务修改了该行数据。
- `幻读(phantom-reads)`: 可以视为 `不可重复读` 的进阶版，A事务中对某一个范围查询的数据两次结果不同。

光看问题其实我们心里应该就有些解法了吧？解决这些问题的常见手法就是加锁：共享锁，排他锁，然后可以做行锁也可以做范围锁，正常我们只要读时加共享锁，写时排他锁，就可以避免以上所有问题。
但是在工程中我们还是需要考虑如何做取舍，即能允许一定的并发，又可以保证安全，Mysql的InnoDB使用MVCC版本控制来解决，具体资料参考：[深入理解MySQL中事务隔离级别的实现原理](https://segmentfault.com/a/1190000025156465)



## Mysql 的Binlog同步
阿里内部大量使用了 Binlog 同步来解决分布式事务，还贡献了一个开源方案 [Canal](https://github.com/alibaba/canal)，但是这里有几个问题要注意:
- Mysql同步模式：默认是异步的，性能较高，但是可能存在不一致。此外还有半同步和同步两种方式，半同步指配置好期望的slave节点数，当同步完指定数量之后才会返回成功，期望slave节点数=all slave节点数时则是同步方式。
- 事务一致：mysql对于binlog的落盘可能是分批的，如果大于binlog的缓冲区则会先行flush，此时就会产生一个event并 [刷盘](https://dev.mysql.com/doc/refman/8.0/en/replication-options-binary-log.html#option_mysqld_binlog-row-event-max-size)，而 canal 也是基于 event来实现消息投递的，但是它在内部实现了一个 [环形缓冲区](https://cloud.tencent.com/developer/article/1661530)，将一个事务的event都汇聚成一个消息投递。
- 顺序：这里顺序包含两个方面，一是写入方写入到Partition/Queue时是否有序，这和是否有多个 P/Q，写入端是否使用并行写入有关；二是消费顺序，这和消费端是否使用并行消费有关。
- 内容：mysql binlog有三个模式，row, statement, mixed，row模式会记录产生变化的数据行的列，作为 where 条件来生成语句，避免产生不一致的运行结果，选取的列取决于 `binlog_row_image` 参数的设置(full, minimal:仅变化列, noblob:除开 text 和 blob的列);statement 记录的是SQL的变化。如果你使用了是 statement 的模式，要注意同步的结果是否满足预期，如果你使用的是 row 模式，要注意日志的传输和存储成本。

## 参考文档
- [Wiki NoSQL](https://zh.wikipedia.org/wiki/NoSQL)
- [一篇文章讲透MySQL为什么要用B+树实现索引](https://cloud.tencent.com/developer/article/1543335)
- [cap](https://cloud.tencent.com/developer/article/1437772): 实际应用中的可用性和CAP可用性并不相同。你应用的可用性多数是通过SLA来衡量的（比如99.9%的正确的请求一定要在一秒钟之内返回成功），但是一个系统无论是否满足CAP可用性其实都可以满足这样的SLA。