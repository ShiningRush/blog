+++
title = "接口安全漫谈"
date = "2021-02-12T14:32:31+08:00"
author = "vincixu"
authorTwitter = "" #do not include @
cover = ""
tags = ["api", "security"]
keywords = ["api", "security"]
description = ""
showFullContent = false
+++

# 接口安全漫谈
针对接口发起的攻击有以下：
- DDos洪水攻击
- 重放攻击
- 数据劫持以及篡改
- 数据监听

## 重放攻击
> 指攻击使用相同的数据包请求发起请求，如果服务端处理不当可能会造成服务宕机，数据被污染等等。

解决方案：在 `QueryString` 添加 `TimeStamp` 与 `Nonce` 字段用于标记请求。
- 放入 `QueryString` 是为了避免添加 CORS 头部
- `TimeStamp` 使用 `Unix` 的秒级时间戳，当服务端检测时间晚于服务端时间时，则拒绝该请求。拒绝阀值根据业务情况自行选择，一般考虑 `30m ~ 1h` 足够。
- `Nonce` 用于在攻击者重放有效时间内的请求时进行区分，如果在有效时间内已经处理了相同的：`路径`、`时间戳`、`随机数`的请求，那么拒绝该重复请求。

## 数据劫持以及篡改
数据劫持和篡改一般出现于中间人攻击，就算现在都是 https 连接可以防止过程中被劫持，但是也防不住攻击者直接从 Client 端修改数据。
防止篡改的可以使用带密钥的签名算法对请求参数进行 Hash 后添加到 Header 中携带过来。


## 加/解密、签名、摘要

- `加/解密` 的过程可逆，用于保护敏感信息在传输过程中被其他人窃听
- `签名` 和 `摘要` 都不可逆，用于检验数据完整性

`加/解密` 算法可以分为 `对称加密` 与 `非对称加密`，区别在于加密过程与解密工程是否使用同一个密钥。

常见对称加密算法：AES, DES
常见非对称加密算法：RSA，ECC
常见哈希算法：MD5，SHA-1，HAMC, FNV, Murmur, DJB, CRC
参考 [Hash](https://github.com/dgryski/dgohash)

以上算法复杂度与安全度成反比，顺序按照算法复杂度由低到高。
如果需要一個快速的 Hash 算法，目前已知最快：murmur > crs > md4 > md5, prometheus 的 hash算法用的md5
参考 [Fastest hash for non-cryptographic uses?](https://stackoverflow.com/questions/3665247/fastest-hash-for-non-cryptographic-uses)

### RSA 密钥格式
RSA 密钥格式有两种：pem, der
内容格式有也有两种：PKCS#1, PKCS#8

pem 格式如下
```
-----BEGIN RSA PUBLIC KEY-----
BASE64 ENCODED DATA
-----END RSA PUBLIC KEY-----
```
der 格式则是对将 pem 格式中的 base64 反解之后的原生文件，它是一个二进制内容的，不方便阅读，如果你按照 `ASN.1 ` 规范解开后，它的内容格式如：
```
RSAPublicKey ::= SEQUENCE {
    modulus           INTEGER,  -- n
    publicExponent    INTEGER   -- e
}
```

PKCS#1 是 RSA 专用格式（上面例子）， PKCS#8 是通用格式( 可用于 ECC )，格式如下：
```
-----BEGIN PUBLIC KEY-----
BASE64 ENCODED DATA
-----END PUBLIC KEY-----
```
 
发现其实只是少了 `RSA` 字样，其他是类似的，但是内容中携带了加密的算法，格式如下：
 ```
 PublicKeyInfo ::= SEQUENCE {
  algorithm       AlgorithmIdentifier,
  PublicKey       BIT STRING
}

AlgorithmIdentifier ::= SEQUENCE {
  algorithm       OBJECT IDENTIFIER,
  parameters      ANY DEFINED BY algorithm OPTIONAL
}
 ```