+++
title = "k8s笔记"
date = "2021-02-12T14:42:54+08:00"
author = ""
authorTwitter = "" #do not include @
cover = ""
tags = ["k8s"]
keywords = ["k8s"]
description = ""
showFullContent = false
+++

# k8s 笔记

这里记录一些常用但容易忘记的内容。

## 存储相关
### ConfigMap 热更新
更新configmap后：
- 使用该configmap的环境变量和命令行参数不会触发更新。
- 使用该configmap挂载的数据卷会触发更新。
- 如果使用子路径方式挂载的数据卷，将不会触发更新。

### PV(PersistentVolume) 与 PVC(PersistentVolumeClaim)

类似 `Pod` 与 `Node` 的关系，每个 `PVC` 可以与一个 `PV` 绑定。它们之间需要满足以下条件：
- 请求大小相同
- `StorageClass` 相同


但是与 `Pod` 的不同处在于：一个 `PVC` 只能与 `PV` 一一对应（因为要求资源大小相等）。
使用 `PV` 我们让一个 `POD` 的存储得以持久化，通常来说我们都会配合 `StatefulSet` 来使用。当然，你完全可以通过 `PV` 让一个 `Deployment` 持有状态，但 `StatefulSet` 会更适合有状态的场景。
同时 `StatefulSet` 中还提供了 `VolumeClaimTemplates` 这样的字段来方便自动创建 `PVC`，配合 `动态PV` 可以轻松驾驭持久存储的体系。

一些 Tips:
- 预先创建 `PV` 再创建 `PVC` 绑定它们的这种方式被称为 `静态PV`，在 `1.12` 出现了 `StorageClass` 来完成自动创建 `PV` 的过程，被称为 `动态PV`。
- `StorageClass` 是一个 `CRD( 自定义资源 )`, 它包含了如何访问数据源的一些细节，如参数，类型等等，为 `动态PV` 提供了支撑。在之前，这些信息都保存在 `PV` 中。
- `StatefulSet` 删除时不会自动删除 `PVC`，需要用户手动删除，这是为了让用户能够恢复数据。其他资源则不会，如 `Deployment`。


## 杂项
- client-go 中提供了很多有用的库：workerqueue(包含限速和延迟队列，多对多且并发安全), leaderelection, resourcelock(类似分布式锁)
- ServiceAccount 已弃用，ServiceAccountName 是最新项
- 每个 ServiceAccount 创建时都会自动生成一个 ServiceAccoutSecret，并且挂载到使用了 ServiceAccount 的 Pod 中
- 可以使用以下命令查看容器中的网络
```
docker inspect -f {{.State.Pid}}    容器id   # 获取容器的pid
nsenter -n -t pid   # 进入容器网络空间
```
简写为
```
nsenter -n -t $(docker inspect -f {{.State.Pid}} dockerid)
```
- ssh 在登录过程中会清理环境，导致 k8s 设置的环境变量都被删除，解决方案可以参考 [](https://stackoverflow.com/questions/34630571/docker-env-variables-not-set-while-log-via-shell)


## 网络模型
k8s 对网络提出了三个基本要求
- pods on a node can communicate with all pods on all nodes without NAT ( 节点上的 POD 不需要通过 NAT 就可以与其他任意节点的 POD 通信 )
- agents on a node (e.g. system daemons, kubelet) can communicate with all pods on that node (节点上的代理程序可以与该节点上的任意 POD 通信 )

下面的这个要求仅针对支持 `HostNetwork` 模式的平台，如 linux
- pods in the host network of a node can communicate with all pods on all nodes without NAT ( 节点上的 hostnetwork POD 不需要通过 NAT 就可以与其他任意节点的 POD 通信 )

实现了这些要求的技术方案有很多，可以参考 [集群网络系统](https://kubernetes.io/zh/docs/concepts/cluster-administration/networking/)，这里列举几个比较知名的开源方案：
- Flannel
- Calico

各大云厂商也都实现这些网络要求以提供 k8s 集群的网络资源需求，由于云厂商在提供容器服务之前基本都有了自己的虚拟化技术而且基本都类似，所以他们都是基于各自的 VPC 网络模型来实现 k8s 的网络方案，比如：
- Google: Google Compute Engine (GCE), 详细情况未知
![Aliyun](/image/k8s-notes/gce.jpg)
- AWS: AWS VPC CNI, 类似 Tencent VPC CNI
![Aliyun](/image/k8s-notes/aws.jpg)
- Azure: Azure CNI, 详细情况未知
- Aliyun: Terway，类似 Tencent VPC CNI，兼容  Flannel
![Aliyun](/image/k8s-notes/ali.jpg)
- Tencent:
  + Global Router: 使用 Overay 的技术建立的容器网路，母机 和 POD 网络互通，但是由于 Overlay 建立在 VPC 下，因此容器不能访问到不同 VPC 下的 POD。
  ![vpc-cni](/image/k8s-notes/gr.jpg)
  + VPC CNI: POD 和 母机 在一个网络平面，这个模式下 POD 和 母机以及 不同 VPC 下的母机都是互通的。
  ![vpc-cni-new](/image/k8s-notes/vpc-cni.jpg)
  + VPC CNI 独立网卡模式: VPC CNI 母机上的 POD 都共用母机的 ENI，新一代的 VPC-CNI 直接将弹性网卡绑定到了 POD 网络命名空间，让其独享 ENi，在实践过程中来看，这个模式提高了网络稳定性，老模型在大包量的情况会出现丢包。
  ![vpc-cni-new](/image/k8s-notes/vpc-cni-new.jpg)

## DNS解析
FQDN(Fully Qualified Domain Name) 代表完全限定域名，在进行域名解析时将不会在后面追加顶级域名，FQDN( Partially Qualified Domain Name ) 部分限定域名，可以根据场景在后方追加顶级域名。
最明显的区别在于 FQDN 后跟一个 `.`, 如 `www.baidu.com.`
- 解析域名时，遵从 `/etc/resolve.conf` 指定的参数, `nameserver` 指代下一级查询的 dns 服务器， `search` 指定在域名为 PQDN 时，尝试追加的基本域名
- k8s 在网络容器初始化时会根据 yaml 的内容来决定如何初始化 pod 的内容，参考 [Pod 与 Service 的 DNS](https://kubernetes.io/zh/docs/concepts/services-networking/dns-pod-service/)
- 在普通 linux 的 `resolve.conf` 配置 `search` 后，如果请求域名中没有 `.` 那么认为是主机名，优先追加 `search` 的域名去查找，最后才查找主机名，比如 `host`，如果含有点，且大于 `options.ndots:n` 当中配置的数量 `n`，则认为是 FQDN 优先查找自身 如 `host.name`，可以通过 `host -a `命令查看解析域名的过程

## CPU 管理策略
k8s 在 linux 上管理容器CPU资源的策略有两种：
- CFS(Completely Fair Scheduler)：使用CFS调度器控制CPU能够使用的额度
- CPUSet: 控制容器能够使用的CPU核

这两种策略通常会结合起来一起生效，比如设置某容器的CFS 的 quota 和 period 比值为 2核，但是 cpu set只绑定了 1 号核，那么容器能够使用的最大CPU也只有 1 核，反之也如此。
k8s 简单来说，当 limit = request 时会为容器绑定专用核，其他情况会绑定共享CPU核。
详细内容可以查看[控制节点上的 CPU 管理策略](https://kubernetes.io/zh/docs/tasks/administer-cluster/cpu-management-policies/)

## k8s 的回收策略
这里涉及到两个方向：
- 已经执行完的历史POD数据需要从ETCD中清理，此部分内容的相关配置在 scheduler 中
- 已经执行完的历史POD数据需要从母机上清理，包括容器记录和磁盘占用，这部分内容在 kubelet 中配置，可以参考 [容器镜像的垃圾收集](https://kubernetes.io/zh/docs/concepts/cluster-administration/kubelet-garbage-collection/)

第一点会影响集群的 APIServer 的性能，进而影响到整个集群，第二点如果不设置的话，有可能会导致母机磁盘写满，进而导致母机 NotReady。通常 kubelet 会在母机磁盘到达一定阈值(默认85%)后自动清理不使用的镜像，但是不会清理历史容器，所以也需要关注全局保留的旧容器数量，默认不限制。


## 健康检查
在 k8s 中容器健康检查分为 两类：
- 存活探测( liveness probe )
- 就绪探测( readyness probe )

这两者的配置都一样，但是产生的结果不同：
- 存活探测：初始状态为满足，防止初始化时直接重启容器，探测失败直接重启容器。
- 就绪探测：初始状态为不满足，防止容器尚未就绪就将流量转发到 pod 中，探测失败会从 Endpoint 中移除。

从上面的角度来考虑，配置健康检查的最佳实践：
- 将存活探测后置于第一轮就绪探测来执行，防止第一轮就绪检查尚未成功时就重启容器
- 就绪探测比存活探测更密集，可以及时发现应用阻塞的情况，再配合存活探测来重启容器


## 平滑更新
平滑更新在 k8s 中其实是相对容易实现的，只要注意以下几点：
- 如果使用了集群外的服务发现机制，比如 consul, zookeeper，这里分情况：
  + 如果没有做缓存，那么只需要考虑做好应用的优雅退出即可
  + 如果做了缓存，那么需要做好端点的健康检查，防止重用旧的端点
- 如果使用 `Service` 来做服务发现，要注意 endpoint 和 kube-proxy 都是异步的，所以保险的方法是使用 `preStop` 的生命周期钩子来等待更新
```
lifecycle:
  preStop:
    exec:
    command: ["/bin/sleep", "15"]
terminationGracePeriodSeconds: 30
```
- 如果使用长连接，请注意在 k8s 1.14 版本以前会存在问题，因为 kube-proxy 会直接删除老规则，导致旧的连接收发包时会直接被丢弃，而在新版本中添加关于 Service 的优雅退出：
  + 首先将老的规则权重更新为 0，这样新建立的连接 (connect) 并不会使用它
  + 老的连接进行收发包时会正常进行
  + pod 被销毁后，也会剔除老规则
- 应用的优雅退出不应该在接收到 SIGTERM 信号立即关闭监听端口，而应该先进行收尾工作：处理当前请求，关闭持有的长连接、Worker等，最后再关闭监听端口，以最大限度容忍由于某些异常原因依然抵达了该服务的请求。

## 以一个 pod 的创建来观察 k8s 的组件协作
一个 POD 要跑起来会经历以下流程
- 请求 k8s APIServer，将资源文件添加到集群，这里会涉及到 kubeConfig 的 认证，授权，以及准入控制
- 当资源被 APIServer 接纳后，`Scheduler` 将会 Watch 到新的资源的产生，并尝试将 Pod 调度到Node上，这里调度涉及几个过程：
  + 预选：将会根据 Pod 的所需资源（包括计算资源、自定义资源、数据卷等等）、NodeName、NodeSelector、亲和性和容忍度，以及Node的健康状态，筛选出一批可容纳Pod的Node，这些过程都分布在不同调度策略中
  + 优选：对上个过程中的Node进行打分，打分维度包括：节点的空闲资源，节点的POD数，Pod的亲和性和容忍度等，从中选出得分最高的节点，对节点和Pod进行绑定
- 当节点的 kubelet watch到有pod与自己进行了绑定则开始创建Pod，流程如下：
  + 通过 CNI 创建网络空间，通过 CRI 创建 Sandbox，相继启动 Init容器与业务容器，再根据其需求判断是否需要使用 CSI 挂载数据卷，最近再进行健康检查与就绪检查。

至此整个 pod 则被拉起，这里值得一提的有几个点：
- CRI 创建 Sandbox 的流程在不同CRI下是不同的，比如 runC 体系下，会拉起 pause 容器，然后业务容器共享 pause 容器的网络命名空间，在kata下则会直接拉起一个虚拟机，其他容器则以进程方式共享这个虚拟机。
- docker 并没有实现 CRI 规范，因此为了支持 docker，k8s 维护了 dockershim 来作为转换组件，docker背后其实跑的也是 containerd，而containerd 可以选择 kata 与 runC 两个运行时,两者都符合 OCI(Open Container Initiative) 规范