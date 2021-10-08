+++
title = "APISIX 插件"
date = "2021-02-12T15:30:56+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["apisix"]
keywords = ["apisix"]
description = ""
showFullContent = false
+++

# APISIX 插件
## 注意点
插件的生命周期查看 [这里](https://github.com/apache/incubator-apisix/blob/e6804360d1712d456cfee31b9d6abb8b16cfffba/doc/architecture-design-cn.md)

插件的编写用法查看 [这里](https://github.com/apache/incubator-apisix/blob/master/doc/plugin-develop-cn.md)

值得注意的一点是，它的插件是按照 `Phase` 分层的。
比如有多个优先级不同的插件：`A > B > C`，他们的`Phase`执行顺序是：
`A(rewrite) > B(rewrite) > C(rewrite) > A (access) > B (access) > C (access)`

插件编写完毕后需要在`config.yaml`的`plugins`节点下加入你的插件名，然后调用`/apisix/admin/plugins/reload` 接口重载下插件即可。
如果发现插件没有生效，或者无法找到插件的提示，可以检查错误日志中的内容，apisix加载插件出错时会记录日志。

## 全局插件
在记录这篇笔记时，官网还没有关于全局插件的用法介绍。
调查了下，apisix有`global rules`的对象可以指定全局的插件。

特此记载了下用法

### 准备一个简单的插件
```lua
local core = require("apisix.core")

local plugin_name = "demo-plugin"

local schema = {
    type = "object",
    -- 这里是插件的属性, 详细参数可以参考其他其他插件
    properties = {
        a_fixed_string = { 
            type = "string"
        }, 
    },
}

local _M = {
    version = 0.1,
    priority = 10086,
    name = plugin_name,
    schema = schema,
}

-- 检查插件的定义是否有效
function _M.check_schema(conf)
    local ok, err = core.schema.check(schema, conf)

    if not ok then
        return false, err
    end

    return true
end

-- 插件的rewrite步骤, 步骤的顺序可以参考文档中的链接
function _M.rewrite(conf, ctx)
    -- 通过conf读取预设的参数
    core.log.warn("get a requert, fixed string: " .. conf.a_fixed_string)
    -- 解开下行的注释可以直接返回指定的请求
    -- return 500, { message = "we got a error" }
end

return _M
```

*启用的方法不再累述，参考上面*

### 创建一个GlobalRule对象
```shell
curl -X PUT \
  https://{apisix_listen_address}/apisix/admin/global_rules/1 \
  -H 'Content-Type: application/json' \
  -d '{
   "plugins": {
   		"demo-plugin": {
   			"a_fixed_string": "FixedString"
   		}
   }
}'
```