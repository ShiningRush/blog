+++
title = "CORS 小记"
date = "2021-02-12T14:35:14+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["http", "cors"]
keywords = ["http", "cors"]
description = ""
showFullContent = false
+++

# CORS(Cross-Origin resource sharing) 笔记
CORS是一种为了解决跨域请求而诞生的规范，详细规范参见 [MDN](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Access_control_CORS)

这里只介绍下一些需要注意的点：
- 跨域问题只存在于`客户端有安全策略限制的情况`，典型的如：浏览器，此时安全策略会阻止你进行跨域访问。如果是自己编写的程序，没有安全策略的限制，则不需要考虑跨域，所以`CORS`通常是前端考虑的问题。
- CORS有`简单请求`和`需预检的请求`区分，只有`需预检的请求`才会发起`OPTIONS`，两者的区分参考上面的MDN。
- 如果需要对凭据进行操作时，比如 携带了`Cookie`，此时需要设置`Access-Control-Allow-Credentials: true`，否则请求将会被浏览器拒绝。
- 在`Access-Control-Allow-Credentials: true`时，无法设置字段为通配符`*`。其他情况下任意字段都可以设置为通配符。这是因为允许携带凭据会带来安全隐患，浏览默认会拒绝 Ajax 的跨域请求，通常 CSRF 攻击只能通过 `script` 标签、`form` 表单 或者 有副作用的 `GET` 来完成，但是一旦开启了 `Access-Control-Allow-Credentials: true`，且不限制 `Origin` 等，那么会有更大的风险遭受 CSRF 攻击。
- Access-Control-Max-Age各个浏览器的最大时长不一样，
  - Firefox caps this at 24 hours (86400 seconds).
  - Chromium (prior to v76) caps at 10 minutes (600 seconds).
  - Chromium (starting in v76) caps at 2 hours (7200 seconds).
  - Chromium also specifies a default value of 5 seconds.
  - A value of -1 will disable caching, requiring a preflight OPTIONS check for all calls.
- 当 `Access-Control-Allow-Origin` 不为 `*` 或者 `静态Origin` 时，比如可以根据访问方的 `Origin` 头部进行动态返回时，要注意设置 `Vary: Origin` 这样浏览器才不会缓存上一次的 `CORS` 返回头，也不会引起 `缓存中毒` 的安全隐患，缓存中毒，指攻击利用浏览器根据 url 缓存返回结果的特性，从自己的攻击站点发起 `XSS` 攻击，此时浏览器将会直接使用上一次的跨域请求回复，而通过攻击者的请求。