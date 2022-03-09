+++
title = "单元测试的必要性"
date = "2021-08-18T16:30:04+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["app-design","unit-test"]
keywords = ["单元测试", "单元测试的必要性"]
description = "很多人都在质疑单元测试的必要性，今天我们就来探讨下，单元测试对一个服务或者框架来说，到底有多重要"
showFullContent = false
+++

分享下个人在过去团队里做单测总结的一些经验。

> 单元测试是不是必须的

我觉得多数情况下是必须的，但可以对覆盖率要求别那么严苛。

每一个服务肯定都存在自己一些核心业务逻辑，这些逻辑我觉得必须要单元测试。但是还有部分和基础设施强相关的逻辑，比如对数据库的查询，这些可以适当要求低覆盖率甚至不做。

举一个老生常谈的例子： A 给 B 转账。
这里面有两个关键业务逻辑
- A, B 账户存在且状态正常
- A 账户余额充足

<br />
可能有的同学会问，还有个事务性呢，A 与 B 的更改必须同时生效。
其实这个不属于业务逻辑，业务逻辑只保证同时提交修改后的 A 与 B 的账户，如何让更改原子地落盘是应用层的逻辑。

按上面提出的几个业务逻辑，那么我们可以有一些 case:

| give | want |
| ---- | ---- |
| A 账户不存在 | error1 |
| B 账户不存在 | error2 |
| A 账户存在，状态冻结 | error3 |
| B 账户存在 ，状态冻结| error4 |
| A, B账户存在，状态正常，A账户余额不足 | error5 |
| 所有都正常 | 正常转账 |

单测可以不用关心是如何读取到 A, B 账户的，从 MariaDb or Mongo or Redis 都无所谓，只需要构建出两个符合要求的账户实体就可以了。单测要确保的是：
**在转账操作结束后，得到预期的 `error` or `已修改的账户`**

在上面这个场景里，如何从 `持久层` 读取实体这个操作，我通常不会做单测，因为失效率相对比较高，导致维护这些单测的成本也高。比如数据库结构变更了，数据迁移了，加了分级缓存了。这些操作都可能引起读取相关的单测变化。而其实维护这些操作所带来的 `业务价值` 是比较低的。

从业务角度来说，其实不关心数据存哪里，怎么读，只关心我们转账的时候，账户的状态是不是按照预期地变化。
说了这么多，其实就一句话：`优先保证核心逻辑的单测`。

如果放在具体的场景来说
- **如果是在维护一个业务**: 我会优先保障业务逻辑的单测，做到业务逻辑 100%  单测覆盖。
- **如果是在维护一个框架 or 组件库**: 由于框架和组件库的特殊性，它们的代码全都是 **核心逻辑** 所以尽量做到全覆盖。
*之所以说尽量，是因为一些对第三方基础设施的拓展函数，不太容易单测。*

<br />
>  单元测试谁来做

研发做更好，单测都是在白盒情况下针对细粒度模块的，由于模块定义是由研发完成，那么在修改时自己去保证相关的 case 的正常工作也是他的责任。

不过很多研发同学都缺乏测试的知识，写得的 case 各种边界条件和异常情况都考虑不够，甚至只写正常的 case …… 
其实如果测试同学也熟悉研发的语言的话，由测试来编写 case 那是极好的，这样相当于 `结对编程` 了:joy:
只是项目成本就上来了，相当于培养了一个项目的备份开发，很多时候都是接受不了的。

### 推广单元测试的阻力
其实从工作以来，接触过的多数开发同学都不太爱写单元测试，觉得写单元测太费时间了，我总结了一些导致这个现象的原因：
- 没有区分代码的优先级
- 不关注设计
- 赶进度

#### 没有区分代码的优先级
这一点在业务服务来说尤为重要，业务服务中最有价值的代码 (业务逻辑) 其实占比没有想象中的那么高，很多时候我们业务中的很多代码都是一些 `应用逻辑`，比如 `日志`，`转化请求为业务实体`，`访问数据库`，`读写文件`。它们对于一个完整的服务来说，是必要的。但是对于业务来说，不太关心。就像产品只看重你是否能快速解决bug，而不管你是直接读代码定位问题，还是根据日志推断的。

所以我们如果能够为服务中的代码进行价值分层，然后优先保障最高级别的单测，这样这项工作的负担就没这么重了。但是有个大前提是：我们的设计要支持我们对其依赖部分的内容进行打桩，这个话题请看下面提到的第二点。

#### 不关注设计
我观察到一个很奇怪的现象，很多研发同学都比较倾向于代码的 `性能` 而不太关注 `设计`。
就像我厂入职面试时常会问一些算法题，却鲜有人会问设计题。
你可能遇到面试官问：
"给你一台4core8G的机器，磁盘上有一个100G大的文件，每行都是一个不超过100长度的字符，请问如何对这个文件进行排序？"
"假设有 100 层的高楼，给你两个完全一样的鸡蛋。请你设计一种方法，能够试出来从第几层楼开始往下扔鸡蛋，鸡蛋会碎。 当然，这个问题还有推广版本，有兴趣的同学可以思考一下。 假设有 n 层楼，给你 k 个完全一样的鸡蛋，请问最坏情况下，至少需要试验多少次才能知道从第几层楼开始往下扔鸡蛋，鸡蛋会碎。"

却很少有面试官问：
"假如地铁起步价是5块钱，10公里内不加价，10公里以上每公里 1.5 元，20 公里以上每公里 2 元，30 公里以上每公里 3 元。每逢节假日时，运营可能准备一些活动，比如圣诞节给总计费进行8折优惠。假设地铁站已经给出且固定，请画出一个UML草图，要求尽量符合SOLID原则，比如以后新增40以上每公里 4 元的计费规则，设计或者代码是如何变化？"

由于这样的大环境倾向，导致很多研发同学只关心程序性能高不高，稳不稳定，却没注意过可维护性好不好。很多时候我接手维护的项目，加了个需求，想加单元测试，却无奈地发现业务代码里面直接强依赖了数据库访问层的代码，根本没法打桩，只能被迫做 `集成测试`。这样的设计也让很多想做单测的同学打了退堂鼓。

#### 赶进度
有的时候为了赶进度，大家不会考虑太多，所以代码都是线性下来，没有切面，没用复用，`ctrl c+v` 用的飞起，很快有了一个可用的项目。
短期来看进度确实提升了，但等到出 bug 和 需要二次开发时，你就明白，什么叫 `技术债` 了。
真是
**裸写一时爽，维护全家。。。**

### 写在最后的一些话
- 一个服务要有单元测试，可以不多，但是起码覆盖核心逻辑。
- 单元测试不但可以保证逻辑的正确，同时也是检验你设计的一项指标。
- 单元测试有时候比文档更容易阅读，很多时候我在阅读一些开源框架的源码时会直接通过函数的单元测试来了解它的主要功能，直接又简洁。
- 一个不可测试的设计，随着复杂度增高，总有一天会崩塌的。
- 设计 和 算法 对工程师而言都是很重要的基本素养，不应该偏重任何一方。


另外毛遂推荐阅读另一篇文章 [如何优雅地在 go 项目中落地测试]() 有兴趣的同学自取。