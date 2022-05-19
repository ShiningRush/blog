+++
title = "领域驱动设计(DDD)入门实践"
date = "2022-01-21T16:05:11+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["", ""]
keywords = ["", ""]
description = ""
showFullContent = false
readingTime = false
+++

# 领域驱动设计(DDD)入门实践
## 背景
我从 2016 年接触到 DDD 开始，到目前( 2022年 )为止，在各种项目中实践过 DDD，包括 医疗设备系统、用户社区、容器云，在这些系统中或深或浅都是用了 DDD 的一些概念模型，因此对 DDD 虽然谈不上熟练，也算是积累了一些宝贵的经验。
本文的出发点：
- 这些年利用 DDD 在解决软件复杂性上所获得的经验，希望能够帮助到更多小伙伴
- 国内互联网近些年都热衷于聊架构、算法、性能，却鲜有在工作中关注代码设计的氛围，希望通过本文能够让读者意识到到代码设计的价值所在。


## 为什么是DDD
我们先来回答这个问题，为什么本文不讲设计模式，不讲重构，不讲架构风格，而偏偏是 DDD。
答案是：`实用性`，为什么说实用呢，我们先来回顾下设计上最广为人知的 `设计模式`。出身的同学想必早就在学校学习过各种设计模式，策略模式，工厂模式，责任链模式等等，但是回顾一下你们的工作生涯中，你们有多少场景用到了这些模式，在整个项目中又占到了多少比例呢，如果你细细回顾下，我相信你使用了设计模式的代码内容在你的项目比重中一定很低。
究其原因在于，设计模式是面向局部场景的经验集合，而非系统性的方法论，在我了解的范围内，DDD 是最全面的设计理论，它能够覆盖到软件开发过程的大多数问题，而其他的一些设计理论大多都是面向特定问题的。


## DDD是什么
领域驱动设计(Domain Driven Design)是一种设计理论，可以先看看 Wiki 的定义：
> Domain-driven design (DDD) is a software design approach focusing on modelling software to match a domain according to input from that domain's experts.
> 
> One concept is that the structure and language of software code (class names, class methods, class variables) should match the business domain. For example, if a software processes loan applications, it might have classes like LoanApplication and Customer, and methods such as AcceptOffer and Withdraw.
> DDD connects the implementation to an evolving model.
> Domain-driven design is predicated on the following goals:
> - placing the project's primary focus on the core domain and domain logic;
> - basing complex designs on a model of the domain;
> - initiating a creative collaboration between technical and domain experts to iteratively refine a conceptual model that addresses particular domain problems.

