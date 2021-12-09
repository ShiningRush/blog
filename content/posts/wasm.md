+++
title = "WASM小记"
date = "2021-02-12T15:15:02+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["prometheus"]
keywords = ["prometheus"]
description = ""
showFullContent = false
+++

# WASM小记
这几天工作需要，了解了一下WASM(WebAssembly)，做个笔记备忘一下

## 背景
WASM 最初诞生的目标是为了解决在浏览器中执行效率的问题，由于JS是一个动态语言，因此对语言其进行优化时会出现一些意向不到的情况，比如上一秒是Object的变量，下一秒就变成了Array，这就导致编译器必须不停地去重新编译成字节码并优化。
为了解决了这个问题，最早出现的解决方案就是 `asm.js`，类似 WASM，它的目标是将 js 编译成一个相对静态的语言，但是本质还是 js，还是需要在浏览器中编译成字节码，可以参考如下例子：
```javascript
function asmJs() {
    'use asm';
    
    let myInt = 0 | 0;
    let myDouble = +1.1;
}
```
为了彻底解决这个问题且获得更好的扩展性，2015年，WASM 诞生了：作为一个规范，它指导各种语言编译成特定格式的字节码，类似：
```
20 00
42 00
.....
```
再借由 `runtime` 去执行它。这种模式让我想到了 java/.net，也都是先编译成中间语言(Intermediate Language)，再到 CLR/JVM 上去执行，如果 .net/java 当年也能支持不同语言编译到IL，是不是就没WASM什么事了？
开个玩笑，.net 也好 java 也好，都是有历史包袱的，而且出发点也并不是为了适配多语言。WASM诞生的目的就是为了提供一个沙盒环境去运行不同的代码，因此在规范上考虑会周全很多。

## 简单介绍
说完了背景，我们来看看 WASM 的工作流程是怎样的。
简单来说，我们构建一个 WASM 模块只需要几个步骤：
1. 通过特定工具，比如 [AssemblyScript](https://github.com/AssemblyScript/assemblyscript),[emscripten](https://github.com/emscripten-core/emscripten)，前者是将 TypeScripe 编译成 WASM 字节码，后者支持多种语言。
2. 找一个运行时来运行 WASM 字节码，现在主流浏览器基本都已经内置了 WASM 的运行时，可以直接执行 WASM 文件

就是这么简单两步，我们就可以体验 WASM 的快乐了，更多的编译和运行工具参考：[Benchmark of WebAssembly runtimes - 2021 Q1](https://00f.net/2021/02/22/webassembly-runtimes-benchmarks/)

## WASI
在 2019年 W3C 将 WASI 纳入规范后，意味着 WASM 已经取得了阶段性的进展，因此他们将目光放到了远方：既然我们已经可以支持多语言在计算层面的编译实现，我们是不是可以看的更远一点，在规范中支持IO（磁盘IO，网络IO），解耦对OS的依赖。
看到这里，又让我想起了另一个技术 dapr(Distributed Application Runtime)，这是微软开源的一个基于云原生的技术，意在解耦程序所用到的所有组件，比如将存储和网络访问都抽象为对 Sidecar 的访问，再由不同的实现去执行具体的逻辑，比如可以装上一个 mysql 的 sidecar，这样存储请求就自然地落入了 mysql，这样程序只需要关心dapr定义的存储规范，而不需要关心具体实现。
dapr 和 wasm 都有一个目的：解耦依赖，不同的是两者解法不同。前者通过规范约束了程序的行为，后者从底层无感知地切换了实现。可以说一个是在应用层的实现，一个是在底层的实现，如果把这两个技术结合在一起，是不是有可能实现极致的弹性伸缩呢，这也许是 Faas 方案的一个出路。

## 总结
虽然 WASM 诞生于 Web，但由于良好的多语言跨平台特性，其实反而在 Web 以外的领域得到的很多实践：
- Istio 方案在 v1.5 后利用了 wasm 实现对 envoy 的动态扩展，同时替代了 mixser 的作用
- ApacheAPISIX 中利用 WASM 来实现多语言插件的执行
- 云平台基于 WASM 来做作为容器运行时，实现毫秒级别的冷启动

同时在 Web 领域也有很多偏计算的应用得益于 WASM 进行了优化和移植，比如 WebGL, Unreal 等游戏引擎。

个人角度来看，WASM技术有几个非常适合的场景：
- Web中移植一些其他语言的库，或者对计算的瓶颈进行优化
- 为应用支持插件机制：WASM的隔离性对于容器场景来说或者偏弱了，但是对于插件场景却特别合适

对于云上环境，我觉得除了 Serverless 的场景可能会有变革性的影响外，目前来看WASM应该没法成为 OCI 的首选，原因如下：
- 隔离性偏弱
- 功能缺失

虽然 docker 的总裁表示，如果 2008 年有 WASM，就不需要 Docker了，WASM 就是这么重要。但我觉得，就算没有 Docker，也一定会有 k8s，会有 CRI，OCI 规范，WASM 最理想也就是一个替换 runC 的选择，而还是会有场景需要隔离性更好的 kata.
可以参考下 [WebAssembly 与 Kubernetes双剑合璧](https://zhuanlan.zhihu.com/p/111057726)