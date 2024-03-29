+++
title = "项目中是否应该大量使用ORM"
date = "2021-08-18T16:34:50+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["app-design","orm"]
keywords = ["ORM"]
description = "本文重在分析ORM在程序设计中的收益与风险，回答 ——是否应该使用ORM—— 这个常见问题"
showFullContent = false
+++

作为一个OO的热爱者，过往工作中使用DDD落地过金融、医疗以及推送类的一些系统，分享下ORM在项目中利弊，不考虑语言。

话题从以下几个角度展开
- 什么是ORM
- 为什么需要ORM
- ORM的优点
- ORM的弊端
- 小结

## 什么是ORM
为了阅读中不产生歧义，想先统一一下认识，我理解的ORM(Object Relational Mapping)是什么。
ORM，它原义指代一种从数据到模型的映射动作，后来泛指进行这类动作的框架或组件。
有的ORM轻量，仅仅处理原始数据到模型的转换，比如.net的[dapper](https://github.com/StackExchange/Dapper)，java的[mybatis](https://github.com/mybatis/mybatis-3)，go的[sqlx](https://github.com/jmoiron/sqlx)，还有非关系数据库的client一般也都会提供相关的操作；
还有的比较重，它可能会帮你管理数据库的连接、事务甚至是关联的对象([entityframework](https://github.com/aspnet/EntityFrameworkCore), [hibernate](https://github.com/hibernate/hibernate-orm), [gorm](https://github.com/jinzhu/gorm))。
它们都属于ORM的范畴，基于此概念，我们继续讨论。

## 为什么需要ORM
在OO还没有这么流行时，初期时人们关注点都在数据上，早期在日本从事.Net开发时，他们会将关系数据库的数据转化为 *DataTable* (它的数据结构就是如其名，有Row有Column)，然后直接显示到界面上，数据流入时也一样。
后来OO兴起了，程序员们发现了建模的好处，纷纷效仿。但问题就来了，很多底层的数据库访问组件返回都是一些元数据，字段名、字段类型等可能都与模型有差异，一个个手写代码去解决这些问题实在是费劲，于是各种ORM框架和组件应运而生。
所以使用ORM的关键点在于 `模型`，我们需要模型时，就需要ORM，而在OOP中几乎到处都是模型的影子，所以你很难避免使用ORM。
那么回到问题来
> 1. 全部使用orm
> 2. orm+sql
> 3. 纯sql

我的理解是，现在的应用你至少要选用一个ORM，可能它很轻量，仅支持SQL，但它至少要帮你反序列化。

## ORM的优点
以下优点不一定所有ORM都有，只是排出来而已。
从上往下越来越重，越重越方便，也越有可能出问题和不够灵活
- 节省手工从元数据映射到模型的成本
- 不用编写SQL
- 不用管理数据库连接
- 不用管理事务
- 支持加载关联对象

## ORM的缺点
以下缺点是使用ORM时可能碰见的，并不代表所有ORM一定存在
1. 元数据映射时类型难以匹配，比如枚举，日期等，大量字段映射时可能会造成性能瓶颈（不同语言的不同ORM支持程度和性能损耗都不一样，仅作参考）
2. 自动编写的SQL可能存在性能缺陷且难以优化
3. 自动的事务管理太弱，面对互联网分布式的需求很难支持
4. 加载关联对象用不好就可能对数据库造成大量查询甚至是慢查询
5. 在一些后台统计页面需要复杂的连表查询，使用模型化的查询无法满足

以上缺点，
- 1的话，大多数生态好的ORM这类问题都比较少，性能问题的话目前为止我没有遇见过服务的瓶颈点是在字段映射的，如果出现了可以具体场景具体分析。
- 245 都可以通过直接执行SQL来避免，这也是为什么大多ORM都带了执行SQL的功能，值得一提的是，如果 5 的场景数据量太大，那建议你使用 数据仓库 或者 数据湖 来完成
- 3的话这个问题可以用一些分布式事务的方案来解决，可以参考我的另一篇文章：[golang 微服务架构中，如何实现分布式事务]()

## 小结
最后小结一下
- 如果你需要构建一个OO应用，那强烈建议你选择一款ORM
- 如果你的应用是一个单机企业应用，不需要考虑分布式，可以考虑选择一款重一点的ORM帮你快速构建应用
- 如果你在构建一个分布式应用，选择一款轻量的ORM，对于带有副作用的请求，你可以使用模型来处理它，如果是普通的查询，你可以使用SQL来搞定。有个专门的模式叫做 [CQRS(命令查询职责分离)](https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/) 就是一种用于解决这类的问题的设计模式，虽然看起来不错，但我一般都不喜欢为此而引入复杂度，直接在应用层写Query就可以了。