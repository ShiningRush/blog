+++
title = "直播与点播"
date = "2023-02-16T21:16:48+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["", ""]
keywords = ["", ""]
description = ""
showFullContent = false
readingTime = false
+++

# 直播与点播
- 直播(Live): 指客户端不能选择观看内容，只能观看服务端的提供内容
- 点播(VOD-VideoOnDemond): 客户端可以选择自己观看的内容

而这其中相关的两个主流协议: 
- RTMP(RealTime Message Protocol): 由 Adobe 创建的，最初的直播协议，支持也最受广泛，它基于 TCP 之上又封装了自己的协议，支持 `推(Publish)/拉(Play)` 模型，延迟低生态好，但是兼容性和网络穿透性差，只支持直播。可以参考[一篇文章搞清楚直播协议RTMP](https://juejin.cn/post/6956240080214327303)
- HLS(HttpLiveSteam): 由 Apple 创建，兼容性和穿透性更好，它由索引文件 m3u8 + segment(ts)文件 组成，由于分片的组成，因此延迟也较大（等于分片大小），支持点播、直播，更多模式可以参考[Example Playlist Files for use with HTTP Live Streaming](https://developer.apple.com/library/archive/technotes/tn2288/_index.html)

正常的mp4格式是不支持流式传输的，需要进行切片参考 [How to output fragmented mp4 with ffmpeg?](https://stackoverflow.com/questions/8616855/how-to-output-fragmented-mp4-with-ffmpeg)

## ffmpeg基础用法
```
ffmpeg
-i 要处理的文件
-o 输出文件名
-c:v 设置编码器，如libvpx-vp9, libxh264等
-c:a 设置音频编码
-crf h264/h265 的参数固定码率因子（CRF），取值 0~51 ，越低越好
-s 输出的分辨率如 1280x720
-r 输出帧率
-d 输出bit率
```

## 音视频基础知识
- 帧率：每秒视频输出的图片数量
- 码率：每秒输出的大小（如果帧率越高，输出的图片越多，在分辨率不变的情况，相应码率肯定会更高）
- 分辨率：视频内容的分辨率，表示由多少像素点组成

正常来说，帧率只要满足 24 fps ，人眼即感知不出来差异度，因此基本都是在控制分辨率和码率，在固定码率的情况下，如果还要加大分辨率，那么编码器只能对原内容进行一些阉割处理，比如色彩信息和打马赛克。

### 透明通道
带透明通道的视频无非就是一组带透明通道的图片(png)，如何制作可以参考 [alpha-transparency-in-chrome-video](https://developer.chrome.com/blog/alpha-transparency-in-chrome-video/)

## 相似度
SSIM(Structural Similarity)，结构相似性，是一种衡量两幅图像相似度的指标，ffmpeg也可以用这种方法来比对视频的相似度，前提是两者具备相同的分辨率和像素格式
> Both video inputs must have the same resolution and pixel format for this filter to work correctly. Also it assumes that both inputs have the same number of frames, which are compared one by one.

像素格式：
- rgb: rgb, argb
- yuv: yuv420, 444, 422


## xavc格式
这是Sony用于记录高清摄影的编码格式，它是遵从h264(5.2)格式做了些优化，以支持高清场景
它的封装格式有两个：
- MXF OP1a: 这是标准封装格式，
- MPEG-4 Part 14: 这是为了在消费者市场推广所使用格式，因为 mp4 格式比较通用

它通常有三个分发标准：
- XAVC: 使用 MXF Op1a 封装，标准格式，里面包含了大多数标准的 XAVC 编解码器
- XAVC S: 使用 MP4 封装，便于分发，编解码同 XAVC
- XAVC HD：使用 MP4封装，里面带了更多的高清编解码器

相关规范:
- xavc的规格地址[这里](https://assets.pro.sony.eu/Web/pdfs/XAVC-technology-pdfs/XAVC_SpecificationOverview_Rev2_2.pdf)
- 更细化的规范，如xavc s, xavc hd在[这里](https://assets.pro.sony.eu/Web/supportcontent/XAVC_Profiles_and_OperatingPoints_210.pdf)