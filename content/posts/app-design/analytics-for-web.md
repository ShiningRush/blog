+++
title = "给你的站点加上统计分析"
date = "2022-01-05T14:59:28+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["app-design","seo"]
keywords = ["hugo","analytics","网站统计","个人网站"]
description = "最近比较好奇自己的个人网站、博客啥的，是不是真的有人看过，平时写的笔记能被搜索到吗、如果搜索不到要考虑做下SEO了，发现光从VPS的服务商那里根本获取不到啥有帮助的信息，所以调查了下目前网站统计接入的方式，记录下如何快速为自己的网站接入统计功能"
showFullContent = false
readingTime = false
+++

# 给你的站点加上统计分析
## 背景
最近比较好奇自己的个人网站、博客啥的，是不是真的有人看过，平时写的笔记能被搜索到吗、如果搜索不到要考虑做下SEO了，发现光从VPS的服务商那里根本获取不到啥有帮助的信息，所以调查了下目前网站统计接入的方式，记录下如何快速为自己的网站接入统计功能

## 统计提供商
虽然市面上有挺多的，但是我建议用两个就行了：
- [百度统计](https://tongji.baidu.com/)
- [谷歌分析](https://analytics.google.com/)

我更喜欢后者的页面，更美观，但美中不足的是，看报告需要翻墙，并且打点的JS文件在国内下载有问题，因此我两个统计都接入了，个人站点免费，因此不用担心费用问题。

### 百度统计
注册完网站后，你可以在最上面，`管理` 的Tab下，查看左边的菜单栏，找到 `代码管理->代码获取`，可以得到一段携带了你站点唯一ID的htlm段，类似
```js
<script>
var _hmt = _hmt || [];
(function() {
  var hm = document.createElement("script");
  hm.src = "https://hm.baidu.com/hm.js?唯一ID";
  var s = document.getElementsByTagName("script")[0]; 
  s.parentNode.insertBefore(hm, s);
})();
</script>
```

把它贴到我们站点的 `head` 内即可，如果你用的 hugo 作为模版引擎，那么可以看看你的主题是否支持拓展 head，你可以看看主题的git仓库下 `themes/yourTheme/layouts/partials` 下是否有类似`extend_xxxx` 的html，如果有，直接在你的网站仓库下添加一个与 `layout` 相同路径的文件，填入以上代码即可。

如果你的主题没有支持这样的扩展性，你也可以直接把它的原 `head.html` 拷贝出来，放在你目录的layouts目录下，保持相同路径，即可替换。

### 谷歌分析
注册账号后，谷歌对于关联域名需要先验证，必须到你的DNS提供商那里解析一条TXT记录，表示这个域名确实归你所有。
然后可以在左边的侧边栏最下面找到 `管理->数据流` 下找到和百度类似的代码片段
```js
<!-- Global site tag (gtag.js) - Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=G-xxxxxxxx"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());

  gtag('config', 'G-xxxxxxxx');
</script>
```