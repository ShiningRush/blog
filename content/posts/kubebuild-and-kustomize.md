+++
title = "Kubebuilder 与 Kustomize"
date = "2021-02-12T14:43:47+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["k8s", "k8s-operator"]
keywords = ["k8s", "k8s-operator"]
description = ""
showFullContent = false
+++

# Kubebuilder 与 Kustomize
## 背景
为什么会把这两者放在一起记录，是因为最近要做一个[Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)，综合考虑了：
- KUDO (Kubernetes Universal Declarative Operator)
- kubebuilder
- Metacontroller
- Operator Framework

以后，决定选用 kuberbuilder作为本次的工具。
大致分析下考虑的因素：
- 因为需要利用 `k8sapi` 所以 `KUOD` 不考虑了，`KUDO` 本身类似于一个控制器，它拓展了 `kubectl` 让你可以使用声明式的方式来发挥 `Operator` 的作用，但本身并不会产生新的 `API`，概念其实很好，可以让使用者尽可能少的接触k8s的代码。
- `Operator FrameWork` 其实和 `kubebuilder` 类似，但它主要在于生态还有与社区的结合，比如`Helm`。就基本的功能而言，`kubebuilder`会更胜一筹，而且这两个项目其实正在考虑互融，参考这里：[What is the difference between kubebuilder and operator-sdk](https://github.com/operator-framework/operator-sdk/issues/1758)
- `Metacontroller`没看，不发表意见


## Kuberbuilder简介
简单的说，其实 `kubebuilder` 就是一个脚手架，把一个 `operator` 需要的东西都准备好了，比如 `CRD`、`Controller`、`RBAC`甚至是 `Webhook`。

我们可以通过它提供的命令行工具 `kubebuilder`，来维护我们的`Operator`，比如初始化、添加`API`、安装到集群等。而它用于维护各种资源清单(Manifests) 的工具就是 [kustomize](https://github.com/kubernetes-sigs/kustomize)

用法参考官方的Book就行了，[点击这里](https://book.kubebuilder.io/)

## Kustomize简介
简单来说，这是一个用于渲染 `k8s` 资源清单的工具。
[点击查看文档](https://github.com/kubernetes-sigs/kustomize/tree/master/docs)