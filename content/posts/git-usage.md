+++
title = "git 小结"
date = "2021-02-12T14:39:40+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["git"]
keywords = ["git"]
description = ""
showFullContent = false
+++
# git 小结
常用命令
- `git remote -v`: 查看绑定了哪些上游
- `git remote add {remoteName} {url}`: 添加上游
- `git clone {url} [{newName}]`: 克隆项目到本地，克隆后的仓库url会默认设置为 `origin` remote
- `git stash push`: 将修改存储暂存区，然后用 `stash pop` 取出，用于在不同分支间同步目录
- `git push [--delete] [--tags] [remote][/branch]`: 推送到目标分支,加上 `--delete` 标记可以删除远端分支, `--tags` 标记表示推送 tag
- `git fetch [remote]`: 拉取远端更改到版本库，但不会同步到工作区
- `git merge [remote][/branch]`: 合并到工作区，会形成一个merge的 commit 记录
- `git rebase [remote][/branch]`: 合并到工作区，不会幸成 merge 记录，更干净
- `git pull [remote]`: 拉取远端变更到版本库，同时 merge 到工作区，相当于 fetch + merge

子模块( Submodule )
适用于想要同时修改两个有依赖仓库的场景。
- `git submodule add [url] [name]`: 添加仓库作为子模块

需要注意的点：
- git checkoutout 和 pull 都会拉取 submodule 绑定的 commit id，如果不指定--recurse-submodules 参数则只拉取但不自动更新，可以执行 git submodule update 更新。
- 要修改 submodule 绑定的版本可以进入子模块目录直接操作，也可以使用submodule update --remote --merge 更新
- 由于 pull checkout 都不会自动更新 submodule，这里存在风险如果开发者没有submodule update 子模块，直接commit -a 的话会将旧版本绑定的submodule commit id推送到新分支，review时需注意这类问题
