+++
title = "AI基础笔记"
date = "2024-02-18T10:44:41+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["", ""]
keywords = ["", ""]
description = ""
showFullContent = false
readingTime = false
+++
# AI基础笔记
## DifussionModel原理
通过前向扩散将图片逐步添加噪声(马尔科夫链)，然后反向使用一个神经网络的参数模型推测回来，这个过程就是训练。
- [Diffusion Models：生成扩散模型](https://yinglinzheng.netlify.app/diffusion-model-tutorial/)


## GPT
通过一系列输入的Token来预测下一次Token的可能性，GPT模型不同之处在于使用Transformer模型来提升对于输入的精确理解：
- G(Generative): 生成式
- P(Pre-trained):预训练，就是通过大量的输入来产出基础模型(无监督学习)，同时还可以通过监督学习(FineTuning) 来微调结果
- T(Transformer): 对输入的预处理，这个是最关键的，也是为什么GPT能理解好自然语言输入的原因
参考；[深度剖析 GPT 的原理、现状与前景](https://sspai.com/post/81036)

## PyTorch 并发推理
当使用一个 pytorch 的模型来进行并发推理时，我们两种实现方式：
- **使用模型提供的BatchInferrence能力**：这个需要将多个请求堆叠后再发送
- **使用多进程来并发推理**：这个需要加载多个模型，导致显存增高，对于AIGC场景来说，单一的SDXL模型可能就会占用完所有显存

