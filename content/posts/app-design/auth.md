+++
title = "API认证授权"
date = "2021-04-13T14:35:06+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["app-design","oauth2"]
keywords = ["oauth2","auth"]
description = "本文记录了常见的API认证(Authentication)和授权(Authorization)协议，包括 ACL, RBAC, ABAC，OAuth2.0"
showFullContent = false
+++

# API认证授权
本文记录了常见的API授权和认证协议。

## 认证( Authentication )与授权( Authorization )
正常的访问中，我们都会涉及到两个阶段：
- 认证：你是谁
- 授权：你能干什么

### 认证
上面已经提到，所谓认证简单来说就是“你是谁”，无论是颁发凭据的一方还是校验的凭据的一方，都是基于这个目的而行动。
常见的认证手段有：
- 密码
- 验证码
- 生物识别（指纹，人脸）
- 证书

而常见的认证协议: 
- Kerberos: 比较完全
- LDAP: 轻量，还需要SASL去完善认证的流量

### 授权
授权是在认证之后识别出你能干什么，这里又分为两类知识：
- 授权模型：你怎么分配权限
- 授权协议：你怎么授予他人使用自己的权限资源

常见的授权模型：
- ACL(Access Control List): 将一组访问权限授予某个账号，类似文件系统的权限
- RBAC(Role Base Access Control): 类似 ACL，但是不再授予给账号，而是角色，最后再将角色和账号绑定，更为灵活。但是缺点也很明显，多了一个角色管理对象。
- ABAC(Attribute Base Access Control): 基于属性的控制，相比 RBAC 的区别在于授予权限的对象变成一类满足特定要求的账号，比如 某个属性大于 X 的账号，给予 Y 权限。该模式最为灵活，但是配置太过复杂。 

常见的授权协议有：
- OAtuh2.0
- OIDC(OpenID Connect)

## OAuth2.0
OAuth2.0 协议是针对授权流程的标准协议，不包括验证，OIDC 是在 OAuth2.0 之上进行的补充，包含了用户身份信息。
有关它的详细内容在 [RFC6749] https://tools.ietf.org/html/rfc6749，同时还有针对撤回Token和自省的 [RFC7009](https://tools.ietf.org/html/rfc7009)、[RFC7662](https://tools.ietf.org/html/rfc7662)。
撤回很容易理解，自省( introspection )指的是验证 token 有效性。
以及为了安全而生的 OAuth2.0 PKCE([RFC7636])。

这里简单介绍下 OAuth2.0 定义的四种模式：
- AuthorizationCode Grant: 授权码模式是最常见也最安全的一种模式，简单来说，接入方( 或者叫应用 ) 通过凭据换取授权码，再通过授权码换取 Token，但是过程中有可能被中间人劫持，所以有了 `PKCE` 来拓展该模式，保证它的安全
- Implict Grant: 隐式模式是授权码模式的简化，去除授权码而直接返回 Token，但是不返回 RefreshToken，这个模式同样也是不安全的
- ResourceOwnerPasswordCredentials: 资源拥有者密钥模式即接入方使用用户本身的密钥来换取凭据，这个一定要接入方是可信的，比如用户的桌面。
- Client Credentials: 这个是接入方以自己的密钥登录的场景。

关于各个模式的详细流程可以查看原生文档，也可以看下 DigitalOcean 的教程，[点击](https://www.digitalocean.com/community/tutorials/an-introduction-to-oauth-2)。
另外有关于 OAuth2.0 除了 PKCE 的授权码模式外都不安全的文章，[查看](https://www.ory.sh/hydra/docs/limitations#resource-owner-password-credentials-grant-type-rocp/)
文章提到了 ROPC 模式是 OAuth为了从 1.0 升到 2.0 但是为了兼容某些 IETF 联盟的大型传统公司而做出的让步， IETF 和 一些熟悉 OAuth 的人都不推荐使用这种模式，这种模式可以用于几种少见的常见：
- 遗留应用向 OAuth 转型
- 无浏览器设备，不过这个部分目前 OAuth 正在起草一个新的流程来补全

关于现在应用的 OAuth 实践，[点击这里](https://www.ory.sh/oauth2-for-mobile-app-spa-browser/)

## Hydra
[hydra](https://github.com/ory/hydra) 是一个实现 OAuth2.0 的开源项目，它抽象了认证规范，将授权和认证实现解耦，非常不错。它属于一个开源组织 `Ory`，它下面还有几个和鉴权相关的项目。
这里简单介绍下一些要点
- 容器默认监听 4444 和 4445 端口，一个 public ，一个 admin，public 用于暴露 OAuth2 协议约定的标准端点，admin 用于管理 Hyrdra，更多细节可以查看 [API文档](https://www.ory.sh/hydra/docs/reference/api)
- 使用前先查看下项目的一些 [限制条件](https://www.ory.sh/hydra/docs/limitations)，最关键的一点应该是要求 Mysql版本必须为 >= 5.7 ，或者使用MariaDB
- 配置可以使用 [配置文件](https://www.ory.sh/hydra/docs/reference/configuration) 或者环境变量，常见命令如下
```
// 命令行参数
-c path/to/config.yaml or --config path/to/config.yaml
--dangerous-force-http // 允许 RedirectURI 为 http，还需要配合 dangerous-allow-insecure-redirect-urls 来使用
--dangerous-allow-insecure-redirect-urls=http://vincixutest.woa.com/auth/accept,http://scrtest.ied.com/auth/accept

// 环境变量
SECRETS_SYSTEM=$SECRETS_SYSTEM -- 配置系统密钥，用于加密
DSN=$DSN -- 数据链接
URLS_SELF_ISSUER=https://localhost:9000/ -- Hydra 部署的Origin，用于签发Issuer
URLS_CONSENT=https://localhost:9001/ -- 实现认证规范的确认授权端点
URLS_LOGIN=https://localhost:9001/ -- 实现认证规范的登录端点
```
- 执行数据库迁移，Hydra 维护了一个数据库的版本管理功能，集成在同一个二进制内，可以使用以下的命令开启数据库迁移
```bash
hydra migrate sql --yes $DSN
```

