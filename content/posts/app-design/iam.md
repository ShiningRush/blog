+++
title = "IAM系统调研"
date = "2022-11-13T17:42:06+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["iam", ""]
keywords = ["iam", ""]
description = ""
showFullContent = false
readingTime = false
+++

# IAM系统调研
记录下最近对业界开源 IAM 系统的调研， IAM: IdentityAccessManagement，简单来说就是用户身份标识与访问管理，通常涉及两个部分：
- **认证(Authentication)**: 通常会涉及到用户登录等操作，用于识别 `你是谁`
- **授权(Authorization)**: 识别身份后就是授权的过程，用于管理 `你能做什么` 

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
- 第三方平台(OAuth2, OIDC)

而常见的认证协议: 
- Kerberos: 比较完全
- LDAP: 轻量，还需要SASL去完善认证的流量

### 授权
授权是在认证之后识别出你能干什么，这里又分为两类知识：
- **授权模型**：你怎么分配权限，通常用于`系统资源`的分配，比如某个用户是否具备某类资源的操作权限。
- **授权协议**：你怎么授予他人使用自己的权限资源，通常用于 `用户资源` 在开放平台等场景使用，

常见的授权模型(复杂度逐步变高)：
- **ACL(Access Control List)**: 将一组访问权限授予某个账号，类似文件系统的权限
- **RBAC(Role Base Access Control)**: 类似 ACL，但是不再授予给账号，而是角色，最后再将角色和账号绑定，更为灵活。但是缺点也很明显，多了一个角色管理对象。
- **Zanzibar**: 谷歌云平台的权限控制模型，建立了一种他们自己的DSL(Object, Subject, Relation, SubjectSet 等)，非常灵活，但是有一定的认知成本。
- **ABAC(Attribute Base Access Control)**: 基于属性的控制，相比 RBAC 的区别在于授予权限的对象变成一类满足特定要求的账号，比如 某个属性大于 X 的账号，给予 Y 权限。该模式最为灵活，但是配置太过复杂，因此业界几乎没用使用。

常见的授权协议有：
- OAtuh2.0
- OIDC(OpenID Connect)

接下来，我们针对几种实际场景来看看这些技术的使用：
- 用户登录: 使用上面提到的`认证`方法。
- 内部用户登录: 
   + SSO: 在统一入口登录后，利用OIDC or OAuth2 协议来获取登录信息
   + 独立登录页: 可以通过SDK也能达到底层数据互通的效果，但是由于登录态维护在各个系统，因此登录态无法互通
- 用户资源授权: 
   + 用户显式授权: 这里通常都是其他系统需要访问特定系统的用户资源，那么一般都是让用户在页面授权，类似QQ, 微信登录
   + 服务账号: 在一些底层系统的权限控制中，系统的用户可能是另一个系统，那么它想控制的是就是自己的系统权限，可以通过`服务账号(ServiceAccount)`的概念，新建服务账号后对服务账号进行业务维度的授权，随后系统使用服务账号进行访问

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


当然最后这两个部分都需要一个动作——鉴权（Introspect），用于检测颁发的凭证是否有效（包括你的身份识别与权限校验）。

下面是找到的一些开源方案：
- [casdoor](https://github.com/casdoor/casdoor)(BSD): 后端为 go，存储为关系数据库(通过xorm支持), 前端为 `react`，项目开始于 2020 下半年，功能相对比较齐全，认证与授权都支持，授权基于Casbin实现，大致看了下项目结构，http框架使用的beego，将 session 信息维护到了数据库中，看起来有点乱，页面也丑了一点。不过它还打通了部分 OOS 用于存储用户的头像信息，算是一个不同之处。
- [logto](https://github.com/logto-io/logto)(MPL): 前后端都是TS, 前端框架为 `react`，开始于 2021 年下半年，如项目名字所示，该项目仅支持 `认证` 而不支持 `授权`，但是界面非常时尚，并且把登录页的调整也考虑了进去，比较专注，但是后端使用 `NodeJS` 的话，不禁对它的性能会有点担忧（虽然说多时候时候服务的性能都取决于开发者，但是Nodejs的单线程模型与解释执行这两点比起具备CSP模型的go来说会有不小的差距。）
- [arkid](https://github.com/longguikeji/arkid)(GPL): 基于python的系统，web也直接通过服务端渲染来实现，项目开始于2019年下半年，感觉没什么特别亮点。
- [Kanidm](https://github.com/kanidm/kanidm)(MPL): 基于 rust 的 IM 系统，和 `logto` 类似，项目开始于 2018 年下半年，但是文档仍然显示 `in progress`，让人不禁怀疑起它的活跃度（虽然代码一直都有人在迁入中）
- [keycloak](https://github.com/keycloak/keycloak)(Apache2.0): 基于 java的IAM方案，web不确定如何构成的，依赖
- [authelia](https://github.com/authelia/authelia)(Apache2.0): 基于 go 的老方案了，从2016年尾开始，活跃度一直都还挺高，可惜也只是一个 IM 方案，backend为 fasthttp, web 为 react，
- [Ory](https://github.com/ory/kratos) 套件：基于 go，包含很多套件，其中 hydra 我有使用过，设计理念非常新颖，Ory套件包含IAM的多数能力，并且各个组件都是相对松耦合的。有意思的是，它的AM组件(keto)是基于 google paper ——Zanzibar（google 用于管理google数百个云上资源，如日历、设备、地图、相片等的通用策略） 实现的。

另外在 casbin 项目下有个叫 [awesome-auth](https://github.com/casbin/awesome-auth) 的项目，记录了很多关于认证与授权的项目，其中有框架类型的、也有项目类型的，这里不展开了。

总的来看，业界虽然有很多开源项目，以及 **IDaas** 这个赛道有几家公司都想参与进来，但是我没有看到真正经过海量数据验证过的方案，虽然 `casdoor`, `logto` 等年轻方案看起来都很不错，但是它们并不是孵化于大型项目中，使用案例也还没有那么充足，让人不禁担心它的可靠性。

这里 Identity 要区分两个很重要的场景：
- **外部**：一般都是产品的用户，通常是一个C端产品才会具备的用户。在开放平台的场景下，也可能是程序、脚本等。
- **内部**：公司内的用户，它可能是一个具体的员工，在开放平台的场景下，也可能是程序、脚本等。一般都是公司内部的系统，比如ERP，运营平台等等。

在不同的场景下，我们又可以进行功能细分：
- 外部：
  + 一键登录: 包括手机验证码、微信、QQ以及其他社交平台的登录
- 内部:
  + SSO: 对于企业内部来说，具备一个SSO对所有员工都是福音
- 内外部通用:
  + 登录: 最原始的密码登录
  + 授权(OAuth2/OIDC): 对于开放平台来说，可能是对内，也可能是对外开放的，这里注意授权和访问控制不同，一个是控制某个用户的权限，而另一个是将用户权限放开，因此AC可能权限粒度会很细，而授权来说，一般粒度会比较粗，通常都是针对产品级别的，比如 `微信-联系-读取` 权限，而且最后都换转化为对 `API` 的控制（甚至一开始就是），一般通过网关来完成，而访问控制通常都是放置在服务中。
  + 访问控制(AccessControl): 内部的SSO，以及开放平台都需要这个能力，相关的实现有ACL(AccessControlList), RBAC, Zanzibar, ABAC等(复杂度由低至高)

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