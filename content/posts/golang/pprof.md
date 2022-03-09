+++
title = "简单易懂的 pprof 手册"
date = "2021-02-12T15:13:52+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["go", "pprof"]
keywords = ["go", "pprof"]
description = ""
showFullContent = false
+++

# 简单易懂的 pprof 手册 
推荐阅读
- [实战Go内存泄露](https://segmentfault.com/a/1190000019222661)
- [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)

完整格式：
```
go tool pprof <format> [options] [binary] <source> ...
```

省略 `<format>` 可以使用交互式的控制台，或者换成为 `-http` 可以启动一个 web 服务器查看性能分析图
```
go tool pprof -http [host]:[port] [options] [binary] <source> ...
```

`<source>` 支持三种格式：
- profile.pb.gz: 压缩的 profile 文件，可以通过 `-proto` 格式输出
- pprof的 http 接口
- legacy_profile： 遗留格式，不推荐

pprof 的 http 服务的路径如 `/debug/pprof/{res}`, `res` 支持以下几种类型采样：
- allocs: 内存分配
- blocks: 阻塞操作
- cmdline: 显示程序启动命令
- goroutine: 协程
- heap: 堆内存信息
- mutex: 锁争用信息
- profile: cpu
- threadcreat: 系统线程
- trace: 程序运行跟踪信息

我们通常可以有几种方式进行采样：
- 在服务运行环境下，使用 `go tool pprof -proto ...` 对运行服务进行采样，并转出为采样结果，再导入回本地分析。这种方法适用于线上生产环境和办公网网络不通时，可以让运维同学帮忙操作。
- 直接使用 `go tool pprof -http ...` 命令在目标环境打开一个web分析台，这种方式实时又直观，适合办公网和目标环境的网络互通的情况。

*allocs 和 heap 采样的信息一致，不过前者是所有对象的内存分配，而 heap 则是活跃对象的内存分配。*

## 采样相关

你可以直接使用 `curl` 命令访问 http 服务获取原始采样信息，如下:
```
curl http://127.0.0.1/debug/pprof/goroutine > goroutine.profile
```
原始采样数据也可以作为 `<source>` 的数据源信息

内存分析时可以指定 `-sample_index`:
- inuse_space: 正在使用，尚未释放的空间
- inuse_object: 正在使用，尚未释放的对象
- alloc_space: 所有分配的空间，包含已释放的
- alloc_objects:  所有分配的对象，包含已释放的

## 常用命令
- 采样堆内存并输出 profile
```
go tool pprof -proto http://127.0.0.1/debug/pprof/heap
```
- 采样堆内存并输出 profile ( 指定路径 )
```
go tool pprof -proto -sample_index=alloc_space -output xxx.pb.gz http://127.0.0.1/debug/pprof/heap
```
- 采样 CPU 30s 并输出 profile
```
go tool pprof -proto http://127.0.0.1/debug/pprof/profile?seconds=5
```
- 对比两个时间点的内存差值
```
go tool pprof -http=:8080 -base xxx.pb.gz xxx.pb.gz
```
- 打开分析文件
```
go tool pprof -http=:8080 ./xxx.pb.gz
```

## 总结
- heap 只能展示堆上面的内存消耗而不能展示栈，所以当 goroutine 泄露时，你无法通过 heap 定位到
- 多用 `-base` 来对比不同时间点的增长 

完整格式如下( go1.14.6 )
```
usage:

Produce output in the specified format.

   pprof <format> [options] [binary] <source> ...

Omit the format to get an interactive shell whose commands can be used
to generate various views of a profile

   pprof [options] [binary] <source> ...

Omit the format and provide the "-http" flag to get an interactive web
interface at the specified host:port that can be used to navigate through
various views of a profile.

   pprof -http [host]:[port] [options] [binary] <source> ...

Details:
  Output formats (select at most one):
    -callgrind       Outputs a graph in callgrind format
    -comments        Output all profile comments
    -disasm          Output assembly listings annotated with samples
    -dot             Outputs a graph in DOT format
    -eog             Visualize graph through eog
    -evince          Visualize graph through evince
    -gif             Outputs a graph image in GIF format
    -gv              Visualize graph through gv
    -kcachegrind     Visualize report in KCachegrind
    -list            Output annotated source for functions matching regexp
    -pdf             Outputs a graph in PDF format
    -peek            Output callers/callees of functions matching regexp
    -png             Outputs a graph image in PNG format
    -proto           Outputs the profile in compressed protobuf format
    -ps              Outputs a graph in PS format
    -raw             Outputs a text representation of the raw profile
    -svg             Outputs a graph in SVG format
    -tags            Outputs all tags in the profile
    -text            Outputs top entries in text form
    -top             Outputs top entries in text form
    -topproto        Outputs top entries in compressed protobuf format
    -traces          Outputs all profile samples in text form
    -tree            Outputs a text rendering of call graph
    -web             Visualize graph through web browser
    -weblist         Display annotated source in a web browser

  Options:
    -call_tree       Create a context-sensitive call tree
    -compact_labels  Show minimal headers
    -divide_by       Ratio to divide all samples before visualization
    -drop_negative   Ignore negative differences
    -edgefraction    Hide edges below <f>*total
    -focus           Restricts to samples going through a node matching regexp
    -hide            Skips nodes matching regexp
    -ignore          Skips paths going through any nodes matching regexp
    -mean            Average sample value over first value (count)
    -nodecount       Max number of nodes to show
    -nodefraction    Hide nodes below <f>*total
    -noinlines       Ignore inlines.
    -normalize       Scales profile based on the base profile.
    -output          Output filename for file-based outputs
    -prune_from      Drops any functions below the matched frame.
    -relative_percentages Show percentages relative to focused subgraph
    -sample_index    Sample value to report (0-based index or name)
    -show            Only show nodes matching regexp
    -show_from       Drops functions above the highest matched frame.
    -source_path     Search path for source files
    -tagfocus        Restricts to samples with tags in range or matched by regexp
    -taghide         Skip tags matching this regexp
    -tagignore       Discard samples with tags in range or matched by regexp
    -tagshow         Only consider tags matching this regexp
    -trim            Honor nodefraction/edgefraction/nodecount defaults
    -trim_path       Path to trim from source paths before search
    -unit            Measurement units to display

  Option groups (only set one per group):
    cumulative
      -cum             Sort entries based on cumulative weight
      -flat            Sort entries based on own weight
    granularity
      -addresses       Aggregate at the address level.
      -filefunctions   Aggregate at the function level.
      -files           Aggregate at the file level.
      -functions       Aggregate at the function level.
      -lines           Aggregate at the source code line level.

  Source options:
    -seconds              Duration for time-based profile collection
    -timeout              Timeout in seconds for profile collection
    -buildid              Override build id for main binary
    -add_comment          Free-form annotation to add to the profile
                          Displayed on some reports or with pprof -comments
    -diff_base source     Source of base profile for comparison
    -base source          Source of base profile for profile subtraction
    profile.pb.gz         Profile in compressed protobuf format
    legacy_profile        Profile in legacy pprof format
    http://host/profile   URL for profile handler to retrieve
    -symbolize=           Controls source of symbol information
      none                  Do not attempt symbolization
      local                 Examine only local binaries
      fastlocal             Only get function names from local binaries
      remote                Do not examine local binaries
      force                 Force re-symbolization
    Binary                  Local path or build id of binary for symbolization
    -tls_cert             TLS client certificate file for fetching profile and symbols
    -tls_key              TLS private key file for fetching profile and symbols
    -tls_ca               TLS CA certs file for fetching profile and symbols

  Misc options:
   -http              Provide web interface at host:port.
                      Host is optional and 'localhost' by default.
                      Port is optional and a randomly available port by default.
   -no_browser        Skip opening a browser for the interactive web UI.
   -tools             Search path for object tools

  Legacy convenience options:
   -inuse_space           Same as -sample_index=inuse_space
   -inuse_objects         Same as -sample_index=inuse_objects
   -alloc_space           Same as -sample_index=alloc_space
   -alloc_objects         Same as -sample_index=alloc_objects
   -total_delay           Same as -sample_index=delay
   -contentions           Same as -sample_index=contentions
   -mean_delay            Same as -mean -sample_index=delay

  Environment Variables:
   PPROF_TMPDIR       Location for saved profiles (default $HOME/pprof)
   PPROF_TOOLS        Search path for object-level tools
   PPROF_BINARY_PATH  Search path for local binary files
                      default: $HOME/pprof/binaries
                      searches $name, $path, $buildid/$name, $path/$buildid
   * On Windows, %USERPROFILE% is used instead of $HOME
```