+++
title = "Android容器化"
date = "2022-12-10T17:31:20+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["", ""]
keywords = ["", ""]
description = ""
showFullContent = false
readingTime = false
+++

# Android容器化
最近调查了下，沉淀一些要点。

开箱即用的开源方案：
- [Anbox](https://github.com/anbox/anbox): 通过网桥将 Android 对硬件的依赖转换为对宿主机 or 软件实现，比如 OpenGL ES 可以用软件模拟渲染，如果host上的显卡支持相关指令，也可以使用宿主机的硬件。- [ReDroid](https://github.com/remote-android/redroid-doc): 类似 Anbox ，用 C 编写了一些核心模块来管理

为什么业界的容器方案都是基于Linux的，哪怕服务器是arm的:
- 既然是容器化方案，肯定不是说直接在服务器上跑个安卓os来运行东西，不然每次都要重装系统，肯定不好
- 业界在linux系统上沉淀了太多了运维工具，这其中有不少是可以直接移植到 arm 架构下的。
