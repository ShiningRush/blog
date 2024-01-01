+++
title = "Python 强类型编程"
date = "2023-12-31T09:33:57+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["python"]
keywords = ["", ""]
description = ""
showFullContent = false
readingTime = false
+++

# Python 强类型编程
python的`dict`类型是我们在开发中很经常使用的，很方便，可以直接对字段寻址
```python
ret = json.load('{"field1": "value"}')
ptrint(ret["field1"])
```

但是这种编码方式也带来了一些弊端：
1. 由于没有定义数据结构，会导致我们只能查看到数据源头来判断数据的构成，维护成本高
2. 对于一些可能为空的字段，都需要单独的编码处理，太过繁琐，导致代码不够强壮

经过一番调查，发现 python 要使用强类型开发，最大的一个卡点就是json的序列化和反序列化，这里记录下调查结果。

## 如何定义数据结构
最常见的方式就是通过构造函数
```python
class User:
    def __init___(self, name: str, age: int):
        self.name = name
        self.age = age
```

但是这个方式显而易见有个问题——由于字段定义和构造函数强绑定，导致字段的增加也会需要修改构造函数的入参和赋值逻辑，很繁琐。
所以python支持一个特性——dataclass，它会自动的生成构造函数和默认值，更贴近其他强类型语言的使用方式。
```python
from dataclasses import dataclass

@dataclass
class User:
    name: str
    age:int = 18
```
如此一来，我们不再需要频繁调整构造函数和逻辑，非常棒。
在使用过程中，又发现了另一个问题——继承。
我们在`User`中定义了一些可选参数（有默认值的）和一些必选参数（无默认值的），这在没有继承的情况是没有问题的，但是一旦继承之后，会出现以下错误。
```python
from dataclasses import dataclass

@dataclass
class User:
    name: str
    age:int = 18

@dataclass
class ChildUser(User):
    address:str
    account: str = "000000"
```
错误：`TypeError: non-default argument 'account' follows default argument`
其原因是因为在继承后，可选参数与必选参数的位置顺序没有得到很好的处理，我们需要加上一个选项`kw_only=True`
```python
from dataclasses import dataclass

@dataclass(kw_only=True)
class User:
    name: str
    age:int = 18

@dataclass
class ChildUser(User):
    address:str
    account: str = "000000"
```

## 强类型的序列化
当结构定义完之后，如果你尝试对这些 class 进行序列化，会得到一个`TypeError: Object of type User is not JSON serializable` 的错误。
其本质原因是因为 `json.load/dumps` 都不支持 `class`，但是 `class` 都包含了一个 `__dict__` 的内置dict可用于操作。
了解以上的一些相关知识后，我们来看看有哪些方法：
- 直接使用内置的dict
```python
u = User()
json.dumps(u.dict)
```
这种方法优点是简单，缺点是适应性太差，如果你的class中还有继承关系，或者是一个数组引用了你的class，都没办法正常序列化

- 添加一个简单的自定义序列化器
```python
class SimpleClassEncoder(JSONEncoder):
    def default(self, o):
        return o.__dict__

if __name__ == "__main__":
    u = User()
    json.dumps(u.dict, cls=SimpleClassEncoder)

    # 反序列化直接使用lambda函数来解决
    body = '{"name": "test"}'
    ret = json.loads(body, object_hook=lambda d: User(**d))
```
这个方法比起之前的直接使用 dict 的适用场景会广很多，但是问题依然是存在的——如果结构体暴露给其他模块使用，则他们不一定知道需要序列化器才能使用

- 添加一个自定义方法 "toJson" or "fromJson" ，这个就不给代码实例了，其实就是在class实现对自己dict的序列化，需要挨个class实现，成本较高
- 使用第三方库——这里看了不少开源库后，发现了 [pydantic](https://github.com/pydantic/pydantic?tab=readme-ov-file)，在很多开源库中都有使用，它除了解决json序列化的问题外，还为强类型编程提供了自定义的校验能力
```python
from pydantic import BaseModel
from typing import List

class User(BaseModel):
    id: int
    name: str = 'John Doe'
    friends: List[int] = []


class SuperUser(User):
    account: str
    pwd: str = "root"

if __name__ == '__main__':
    # demos = [Demo(name="name1"), Demo(name="name2")]
    u = SuperUser(id=1, account="ss")
    print(u.model_dump_json())
```


整体看下来，使用 `pydantic` 是整体看下来最全面的，唯一的缺点是需要引入第三方库。