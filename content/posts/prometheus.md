+++
title = "Prometheus笔记"
date = "2021-02-12T15:15:02+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["prometheus"]
keywords = ["prometheus"]
description = ""
showFullContent = false
+++

# Prometheus笔记
## Prometheus远程存储优化

Prometheus 以每两个小时为一个块，存储到磁盘中，内存保存近两个小时的内容，同时也会写 `WAL(write-ahead-log)`，预写日志文件wal以128MB的段存储在目录中。这些文件包含尚未压缩的原始数据，因此它们比常规的块文件要大。Prometheus将至少保留3个预写日志文件，但是高流量服务器可能会看到三个以上的WAL文件，因为它需要保留至少两个小时的原始数据。

这些块会在后台被压缩成更大的块以节省磁盘，最大长度取决于 `storage.tsdb.max-block-duration` 参数的设置，默认为保留时间的 `10%`。
一般来说 `storage.tsdb.max-block-duration` = `storage.tsdb.min-block-duration` = `2h` 就相当于禁用了压缩功能，
但是别让它们低于 `2h`，否则有以下的问题：
- 落盘过于频繁，这会很大程度影响 `Promethues` 的吞吐量。
- 由于 `WAL` 至少保留两个小时，所以这部分的内存是没办法释放的。

参考[Remote write tuning](https://prometheus.io/docs/practices/remote_write/)

## rate vs irate
它们都用于计算一个时间区间内指标的每秒变化率，两者的决定性差别在于:
- rate 使用整个时间区间所有点计算出平均变化速率，它会削平峰值
- irate 使用时间区间内最后的两个数据点作为变化速率

从它们两者的计算公式我们可以得到以下推论：
- 当选择的计算区间内仅仅包含两个数据点时，rate 和 irate 没有区别
- 我们使用 rate 来查看某个数据在较长一段时间内的变化趋势，它会消除一些细节防止影响趋势的展示
- 使用 irate 来查看某个数据在一段时间内的细小抖动和完整的变化。
- 使用 rate 时尽量选择较长的时间，而 irate 则反之（太长会丢失很多变化）

## relabel_config vs metric_relabel_configs
`relabel_config` 发生在 `metric_relabel_configs`，通常用于挑选 target 

## 如何处理prometheus 的返回结果
promehtues的返回结果是异构的，类似 `[1610589855,  "11.778672386876817"]`，除了定义异构结构体外，看了下 client 的代码，原来使用了如下的方式实现：
```go
// SamplePair pairs a SampleValue with a Timestamp.
type SamplePair struct {
	Timestamp Time
	Value     SampleValue
}

// MarshalJSON implements json.Marshaler.
func (s SamplePair) MarshalJSON() ([]byte, error) {
	t, err := json.Marshal(s.Timestamp)
	if err != nil {
		return nil, err
	}
	v, err := json.Marshal(s.Value)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[%s,%s]", t, v)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SamplePair) UnmarshalJSON(b []byte) error {
	v := [...]json.Unmarshaler{&s.Timestamp, &s.Value}
	return json.Unmarshal(b, &v)
}
```

它将两个字段重组成了一个 `json.Unmarshaler` 数组，然后用它来接受结果集。

## 使用 Promehtues 求取一段时间内的最大值
可以使用 `max_over_time(query_range[range])` ，请求最大值时为了只返回一个点，可以使用 `query` 接口，而不是 `query_range`