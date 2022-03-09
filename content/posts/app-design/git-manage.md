+++
title = "Git的分支管理模型"
date = "2022-01-11T15:40:49+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["app-design",""]
keywords = ["git","branch"]
description = ""
showFullContent = false
readingTime = false
+++

# Git的分支管理模型
Git的分支管理模型大致有三种：
- [GitFlow](https://nvie.com/posts/a-successful-git-branching-model/)
- [GithubFlow](https://docs.github.com/en/get-started/quickstart/github-flow)
- [GitlabFlow](https://docs.gitlab.com/ee/topics/gitlab_flow.html)

这里大致介绍下
## GitFlow
分支分为：
- master
- develop

各类开发分支(feat, chore, fix)等都是合往 develop，稳定的版本发布定期从 develop 合往 master。
这个开发流程十分简单，但是问题也很明显，它是基于版本发布的，不适用当下公司内部快速迭代的节奏。

## GithubFlow
和你在 Github 上面维护代码是相同的方式，各个开发分支都通过 `pr` 合入到 master，然后自动部署到目标环境，这个流程在很多公司都得到了实践。
缺点在于如果你有多个环境要持续部署，那么需要维护多个环境分支，如何在多个分支间保持一致性就成了一个问题。

## GitlabFlow
GitlabFlow 融合了 GitFlow 和 GithubFlow 两种模式，它首先给出了在 GithubFlow 下多环境分支的解法：规定上下游。你永远只能从上游同步到下游。
详细的模型可以参考上文的链接，我这里不展开了，有点多。

## 个人思考
通过上面的一些标准 Flow，其实我们能看出一些问题了，在这些各类的flow上，其实有几点本质上的不同：
- 版本发布：你交付到客户的产品都是版本化的产物，比如客户端产品。这种方式讲究质量控制和版本规划，但是对于公司内部不断迭代的产品来说，很有可能会失控。
- 持续发布：在公司内部不断迭代的产品会很热爱这个方式，比如内部系统（中台，CRM等）

Gitflow 其实只满足了 `版本发布` 的需求，而 GithubFlow 只考虑了持续发布的需求，因此才出现了 GitlabFlow 来结合两者的优势。
但是 GitlabFlow 其实只是一个理论指导，什么意思？就是说缺乏具体的规范，比如从上游合入到下游，到底用 `merge` 还是 `cherry-pick` ？怎么确定一个良好的上下游？

个人思考过后，沉淀一些经验在这里。
首先考虑自己产品形态是怎样的？
- 需不需要持续发布：如果你是一个开源的产品，没有需要持续发布的环境，那么你肯定没有这样的需求。但是如果你是个公司项目，先不考虑生产环境，你的测试环境一定是持续发布会更好。
- 需不需要版本发布：通常来公司内部系统其实不太注重版本，但是如果你的系统举足轻重，那么将生产环境作为版本发布来看会是一个不错的选择。

如果你需要持续发布，那么你`需要`管理 `环境分支(env branch)`，类似 `prod(duction)`, `pre(-production)`, `test`等，这些分支应该只能通过 mr/pr 合入，当合入后则自动发布到对应环境。
如果你需要版本发布，那么你`可能需要`管理 `发布分支(release branch)`，这个可以参考 GitFlow 中的定义，便于以后进行 hotfix。

就个人经验来说，我觉得有以下几点是应该要做到的：
- 发起 mr/pr 时要检测来源分支是否满足规范
- 如果你维护的好
- 如果一个项目的维护人只有一个或很少，GithubFlow 是个不错的选择，可以结合 Tag 来完成 `持续交付`，在容器环境下，其实持续部署带来的收益并没有很高
- 如果一个项目的维护人很多，那么引入环境分支除了可以做到持续部署外还可以有效避免各个维护人在目标环境的冲突，但是也引入了一些问题：怎么在环境分支间同步变更，确保变更不会遗漏，这里建议引入上下游的形式来避免预期外的代码迁入，类似 GitlabFlow，但是形式不太一样：
  + 上游的顺序为 dev, test, pre, prod，master只能由上游往下游merge
  + 特性分支必须由 master 迁出
  + 特性分支开发完毕后，先合入环境分支（不删除特性分支）
  + 测试没有问题后，如果确认当前环境的所有变更都是可以发往下游的环境分支的，直接自动同步过去。否则必须由 特性分支 各自合往各个环境分支。
  + 最后合入到 master 后，删除特性分支。

可以看到我最后提出的工作模式其实不同点在于完善了 GitlabFlow的具体流程，同时补充了如果环境包含了其他人的变更时，应该怎么操作的指引。
