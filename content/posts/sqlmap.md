+++
title = "SqlMap用法"
date = "2022-06-04T11:56:50+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["", ""]
keywords = ["", ""]
description = ""
showFullContent = false
readingTime = false
+++
# SqlMap简易用法
[SqlMap](https://github.com/sqlmapproject/sqlmap)是开源的SQL注入工具，它既可用于探测，也可以用于注入。

首先clone到本地：
```bash
git clone --depth 1 https://github.com/sqlmapproject/sqlmap.git sqlmap-dev
```

确保本地具有 Python，版本：2.6, 2.7 and 3.x on any platform.

SqlMap流程简单来说分为：
1. 探测
2. 注入

探测首先需要确定你的target，target支持多种形式，最常见的形式是url：
```bash
python ./sqlmap.py -u "http://127.0.0.1:8080/sqli?id=1&name=1"
```
如果你需要指定方法，可以添加 `--method` 指定方法，更多参数参考 [sqlmap用户手册](https://sqlmap.kvko.live/usage)

探测完毕后，可以发起注入，注入参数可以参考用户手册，探测和注入的结果都会保存在本地session文件中，以便于复用