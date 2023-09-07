# CentOS7 Install K8s 1.24.1

## 前言

kubeadm 是 Kubernetes 官方提供的用于快速安部署 Kubernetes 集群的工具，伴随 Kubernetes 每个版本的发布都会同步更新，kubeadm 会对集群配置方面的一些实践做调整，通过实验 kubeadm 可以学习到 Kubernetes 官方在集群配置上一些新的最佳实践。

## 准备

### 系统配置

在安装之前，需要先做好如下准备。3 台 CentOS 7.9 主机如下：

```shell
cat <<EOF >> /etc/hosts
10.0.50.10    master
10.0.50.11    worker1
10.0.50.12    worker2
EOF
```

升级内核到4.17版本以上

在各个主机上完成下面的系统配置。如果各个主机启用了防火墙策略，需要开放 Kubernetes 各个组件所需要的端口，可以查看 Ports and Protocols 中的内容, 开放相关端口或者关闭主机的防火墙。

禁用防火墙:

```shell
systemctl disable --now firewalld
```

禁用SELINUX：

```shell
setenforce 0

sed -i s/^SELINUX=enforcing/SELINUX=disabled/ /etc/selinux/config
```

创建 /etc/modules-load.d/containerd.conf 配置文件:

```shell
cat << EOF > /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF
```

执行以下命令使配置生效:

```shell
modprobe overlay
modprobe br_netfilter
```

创建 /etc/sysctl.d/99-kubernetes-cri.conf 配置文件：

```shell
cat << EOF > /etc/sysctl.d/99-kubernetes-cri.conf
net.ipv4.ip_forward = 1
user.max_user_namespaces = 28633
vm.swappiness = 0
EOF
```

执行以下命令使配置生效:

```shell
sysctl -p /etc/sysctl.d/99-kubernetes-cri.conf
```

### 配置服务器支持开启ipvs的前提条件

由于 ipvs 已经加入到了内核的主干，所以为 kube-proxy 开启 ipvs 的前提需要加载以下的内核模块：

```text
ip_vs
ip_vs_rr
ip_vs_wrr
ip_vs_sh
nf_conntrack
```

在各个服务器节点上执行以下脚本:

```shell
cat <<EOF > /etc/sysconfig/modules/ipvs.modules
#!/bin/bash
modprobe -- ip_vs
modprobe -- ip_vs_rr
modprobe -- ip_vs_wrr
modprobe -- ip_vs_sh
modprobe -- nf_conntrack
EOF

chmod 755 /etc/sysconfig/modules/ipvs.modules && bash /etc/sysconfig/modules/ipvs.modules && lsmod | grep -e ip_vs -e nf_conntrack_ipv4
```

上面脚本创建了的 /etc/sysconfig/modules/ipvs.modules 文件，保证在节点重启后能自动加载所需模块。 使用 lsmod | grep -e ip_vs -e nf_conntrack_ipv4 命令查看是否已经正确加载所需的内核模块。接下来还需要确保各个节点上已经安装了 ipset 软件包，为了便于查看 ipvs 的代理规则，最好安装一下管理工具 ipvsadm。

```shell
yum install -y ipset ipvsadm
```

如果不满足以上前提条件，则即使 kube-proxy 的配置开启了 ipvs 模式，也会退回到 iptables 模式。

### 部署容器运行时 Containerd

在各个服务器节点上安装容器运行时 Containerd。下载Containerd的二进制包:

```shell
wget https://github.com/containerd/containerd/releases/download/v1.6.4/cri-containerd-cni-1.6.4-linux-amd64.tar.gz
```

cri-containerd-cni-1.6.4-linux-amd64.tar.gz 压缩包中已经按照官方二进制部署推荐的目录结构布局好。 里面包含了 systemd 配置文件，containerd 以及 cni 的部署文件。 将解压缩到系统的根目录 / 中:

```shell
tar -zxvf cri-containerd-cni-1.6.4-linux-amd64.tar.gz -C /
# etc/
# etc/systemd/
# etc/systemd/system/
# etc/systemd/system/containerd.service
# etc/crictl.yaml
# etc/cni/
# etc/cni/net.d/
# etc/cni/net.d/10-containerd-net.conflist
# usr/
# usr/local/
# usr/local/sbin/
# usr/local/sbin/runc
# usr/local/bin/
# usr/local/bin/critest
# usr/local/bin/containerd-shim
# usr/local/bin/containerd-shim-runc-v1
# usr/local/bin/ctd-decoder
# usr/local/bin/containerd
# usr/local/bin/containerd-shim-runc-v2
# usr/local/bin/containerd-stress
# usr/local/bin/ctr
# usr/local/bin/crictl
# ......
# opt/cni/
# opt/cni/bin/
# opt/cni/bin/bridge
# ......
```

注意经测试 cri-containerd-cni-1.6.4-linux-amd64.tar.gz 包中包含的 runc 在 CentOS 7 下的动态链接有问题，这里从 runc 的 github 上单独下载 runc，并替换上面安装的 containerd 中的 **runc**:

```shell
wget https://github.com/opencontainers/runc/releases/download/v1.1.2/runc.amd64
mv runc.amd64 /usr/local/sbin/runc
chmod +x /usr/local/sbin/runc
```

接下来生成 containerd 的配置文件:

```shell
mkdir -p /etc/containerd
containerd config default > /etc/containerd/config.toml
```

根据文档 Container runtimes 中的内容，对于使用 systemd 作为 init system 的 Linux 的发行版，使用 systemd 作为容器的 cgroup driver 可以确保服务器节点在资源紧张的情况更加稳定，因此这里配置各个节点上 containerd 的 cgroup driver 为 systemd。修改前面生成的配置文件 /etc/containerd/config.toml：

```text
# sed -i s/"SystemdCgroup = false"/"SystemdCgroup = true"/ /etc/containerd/config.toml

[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
  ...
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
    SystemdCgroup = true
```

再修改 /etc/containerd/config.toml 中的

```text
# sed -i s@'sandbox_image = "k8s.gcr.io/pause:3.6"'@'sandbox_image = "registry.aliyuncs.com/google_containers/pause:3.7"'@ /etc/containerd/config.toml

[plugins."io.containerd.grpc.v1.cri"]
  ...
  sandbox_image = "registry.aliyuncs.com/google_containers/pause:3.7"
```

修改 /etc/containerd/config.toml的国内镜像代理

```text
[plugins."io.containerd.grpc.v1.cri".registry.mirrors]

[plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
  endpoint = ["http://hub-mirror.c.163.com", "https://docker.mirrors.ustc.edu.cn"]
```

配置 containerd 开机启动，并启动 containerd

```shell
systemctl enable containerd --now
```

使用 crictl 测试一下，确保可以打印出版本信息并且没有错误信息输出:

```shell
crictl version
# Version:  0.1.0
# RuntimeName:  containerd
# RuntimeVersion:  v1.6.4
# RuntimeApiVersion:  v1alpha2
```

## 使用 kubeadm 部署 Kubernetes

### 安装 kubeadm 和 kubelet

下面在各节点安装 kubeadm ,kubectl和 kubelet：

```shell
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=0
gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg
        http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF

yum makecache fast
yum install -y kubelet-1.24.1 kubeadm-1.24.1 kubectl-1.24.1
```

运行 kubelet --help 可以看到原来 kubelet 的绝大多数命令行 flag 参数都被 DEPRECATED 了，官方推荐我们使用 --config 指定配置文件，并在配置文件中指定原来这些 flag 所配置的内容。具体内容可以查看这里 Set Kubelet parameters via a config file。最初 Kubernetes 这么做是为了支持动态 Kubelet 配置（Dynamic Kubelet Configuration），但动态 Kubelet 配置特性从 k8s 1.22 中已弃用，并在 1.24 中被移除。如果需要调整集群汇总所有节点 kubelet 的配置，还是推荐使用 ansible 等工具将配置分发到各个节点。kubelet 的配置文件必须是 json 或 yaml 格式，具体可查看这里。Kubernetes 1.8 开始要求关闭系统的 Swap，如果不关闭，默认配置下 kubelet 将无法启动。 关闭系统的 Swap 方法如下:

```shell
sed -i '/ swap / s/^\(.*\)$/# \1/g' /etc/fstab
swapoff -a
free -m 
```

### 使用 kubeadm init 初始化集群

在各节点开机启动 kubelet 服务：

```shell
systemctl enable kubelet.service
```

使用 kubeadm config print init-defaults --component-configs KubeletConfiguration 可以打印集群初始化默认的使用的配置：

```yaml
apiVersion: kubeadm.k8s.io/v1beta3
bootstrapTokens:
- groups:
  - system:bootstrappers:kubeadm:default-node-token
  token: abcdef.0123456789abcdef
  ttl: 24h0m0s
  usages:
  - signing
  - authentication
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: 1.2.3.4
  bindPort: 6443
nodeRegistration:
  criSocket: unix:///var/run/containerd/containerd.sock
  imagePullPolicy: IfNotPresent
  name: node
  taints: null
---
apiServer:
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta3
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns: {}
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: k8s.gcr.io
kind: ClusterConfiguration
kubernetesVersion: 1.24.0
networking:
  dnsDomain: cluster.local
  serviceSubnet: 10.96.0.0/12
scheduler: {}
---
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
  anonymous:
    enabled: false
  webhook:
    cacheTTL: 0s
    enabled: true
  x509:
    clientCAFile: /etc/kubernetes/pki/ca.crt
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 0s
    cacheUnauthorizedTTL: 0s
cgroupDriver: systemd
clusterDNS:
- 10.96.0.10
clusterDomain: cluster.local
cpuManagerReconcilePeriod: 0s
evictionPressureTransitionPeriod: 0s
fileCheckFrequency: 0s
healthzBindAddress: 127.0.0.1
healthzPort: 10248
httpCheckFrequency: 0s
imageMinimumGCAge: 0s
kind: KubeletConfiguration
logging:
  flushFrequency: 0
  options:
    json:
      infoBufferSize: "0"
  verbosity: 0
memorySwap: {}
nodeStatusReportFrequency: 0s
nodeStatusUpdateFrequency: 0s
rotateCertificates: true
runtimeRequestTimeout: 0s
shutdownGracePeriod: 0s
shutdownGracePeriodCriticalPods: 0s
staticPodPath: /etc/kubernetes/manifests
streamingConnectionIdleTimeout: 0s
syncFrequency: 0s
volumeStatsAggPeriod: 0s
```

从默认的配置中可以看到，可以使用 imageRepository 定制在集群初始化时拉取 k8s 所需镜像的地址。基于默认配置定制出本次使用 kubeadm 初始化集群所需的配置文件 kubeadm.yaml：

```yaml
apiVersion: kubeadm.k8s.io/v1beta3
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: 10.0.50.10
  bindPort: 6443
nodeRegistration:
  criSocket: unix:///run/containerd/containerd.sock
  taints:
  - effect: PreferNoSchedule
    key: node-role.kubernetes.io/master
---
apiVersion: kubeadm.k8s.io/v1beta3vir
kind: ClusterConfiguration
kubernetesVersion: v1.24.1
imageRepository: registry.aliyuncs.com/google_containers
networking:
  podSubnet: 10.244.0.0/16
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
cgroupDriver: systemd
failSwapOn: false
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: ipvs
```

这里定制了 imageRepository 为阿里云的 registry，避免因 gcr 被墙，无法直接拉取镜像。criSocket 设置了容器运行时为 containerd。 同时设置 kubelet 的 cgroupDriver 为 systemd，设置 kube-proxy 代理模式为 ipvs。在开始初始化集群之前可以使用 kubeadm config images pull --config kubeadm.yaml 预先在各个服务器节点上拉取所 k8s 需要的容器镜像。

```shell
kubeadm config images pull --config kubeadm.yaml
# [config/images] Pulled registry.aliyuncs.com/google_containers/kube-apiserver:v1.24.1
# [config/images] Pulled registry.aliyuncs.com/google_containers/kube-controller-manager:v1.24.1
# [config/images] Pulled registry.aliyuncs.com/google_containers/kube-scheduler:v1.24.1
# [config/images] Pulled registry.aliyuncs.com/google_containers/kube-proxy:v1.24.1
# [config/images] Pulled registry.aliyuncs.com/google_containers/pause:3.7
# [config/images] Pulled registry.aliyuncs.com/google_containers/etcd:3.5.3-0
# [config/images] Pulled registry.aliyuncs.com/google_containers/coredns:v1.8.6
```

接下来使用 kubeadm 初始化集群，选择 master 作为 Master Node，在 master 上执行下面的命令：

```shell
kubeadm init --config kubeadm.yaml
# [init] Using Kubernetes version: v1.24.1
# [preflight] Running pre-flight checks
# [preflight] Pulling images required for setting up a Kubernetes cluster
# [preflight] This might take a minute or two, depending on the speed of your internet connection
# [preflight] You can also perform this action in beforehand using 'kubeadm config images pull'
# [certs] Using certificateDir folder "/etc/kubernetes/pki"
# [certs] Generating "ca" certificate and key
# [certs] Generating "apiserver" certificate and key
# [certs] apiserver serving cert is signed for DNS names [kubernetes kubernetes.default kubernetes. default.svc kubernetes.default.svc.cluster.local master] and IPs [10.96.0.1 10.0.50.10]
# [certs] Generating "apiserver-kubelet-client" certificate and key
# [certs] Generating "front-proxy-ca" certificate and key
# [certs] Generating "front-proxy-client" certificate and key
# [certs] Generating "etcd/ca" certificate and key
# [certs] Generating "etcd/server" certificate and key
# [certs] etcd/server serving cert is signed for DNS names [localhost master] and IPs [10.0.50.10   127.0.0.1 ::1]
# [certs] Generating "etcd/peer" certificate and key
# [certs] etcd/peer serving cert is signed for DNS names [localhost master] and IPs [10.0.50.10 127.  0.0.1 ::1]
# [certs] Generating "etcd/healthcheck-client" certificate and key
# [certs] Generating "apiserver-etcd-client" certificate and key
# [certs] Generating "sa" key and public key
# [kubeconfig] Using kubeconfig folder "/etc/kubernetes"
# [kubeconfig] Writing "admin.conf" kubeconfig file
# [kubeconfig] Writing "kubelet.conf" kubeconfig file
# [kubeconfig] Writing "controller-manager.conf" kubeconfig file
# [kubeconfig] Writing "scheduler.conf" kubeconfig file
# [kubelet-start] Writing kubelet environment file with flags to file "/var/lib/kubelet/  kubeadm-flags.env"
# [kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
# [kubelet-start] Starting the kubelet
# [control-plane] Using manifest folder "/etc/kubernetes/manifests"
# [control-plane] Creating static Pod manifest for "kube-apiserver"
# [control-plane] Creating static Pod manifest for "kube-controller-manager"
# [control-plane] Creating static Pod manifest for "kube-scheduler"
# [etcd] Creating static Pod manifest for local etcd in "/etc/kubernetes/manifests"
# [wait-control-plane] Waiting for the kubelet to boot up the control plane as static Pods from   directory "/etc/kubernetes/manifests". This can take up to 4m0s
# [apiclient] All control plane components are healthy after 15.502634 seconds
# [upload-config] Storing the configuration used in ConfigMap "kubeadm-config" in the "kube-system"   Namespace
# [kubelet] Creating a ConfigMap "kubelet-config" in namespace kube-system with the configuration   for the kubelets in the cluster
# [upload-certs] Skipping phase. Please see --upload-certs
# [mark-control-plane] Marking the node master as control-plane by adding the labels: [node-role. kubernetes.io/control-plane node.kubernetes.io/exclude-from-external-load-balancers]
# [mark-control-plane] Marking the node master as control-plane by adding the taints [node-role.  kubernetes.io/master:PreferNoSchedule]
# [bootstrap-token] Using token: xflg2x.ehyvpkhdlfzm4n9h
# [bootstrap-token] Configuring bootstrap tokens, cluster-info ConfigMap, RBAC Roles
# [bootstrap-token] Configured RBAC rules to allow Node Bootstrap tokens to get nodes
# [bootstrap-token] Configured RBAC rules to allow Node Bootstrap tokens to post CSRs in order for  nodes to get long term certificate credentials
# [bootstrap-token] Configured RBAC rules to allow the csrapprover controller automatically approve   CSRs from a Node Bootstrap Token
# [bootstrap-token] Configured RBAC rules to allow certificate rotation for all node client   certificates in the cluster
# [bootstrap-token] Creating the "cluster-info" ConfigMap in the "kube-public" namespace
# [kubelet-finalize] Updating "/etc/kubernetes/kubelet.conf" to point to a rotatable kubelet client   certificate and key
# [addons] Applied essential addon: CoreDNS
# [addons] Applied essential addon: kube-proxy
# 
# Your Kubernetes control-plane has initialized successfully!
# 
# To start using your cluster, you need to run the following as a regular user:
# 
  # mkdir -p $HOME/.kube
  # sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  # sudo chown $(id -u):$(id -g) $HOME/.kube/config
# 
# Alternatively, if you are the root user, you can run:
# 
  # export KUBECONFIG=/etc/kubernetes/admin.conf
# 
# You should now deploy a pod network to the cluster.
# Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  # https://kubernetes.io/docs/concepts/cluster-administration/addons/
# 
# Then you can join any number of worker nodes by running the following on each as root:
# 
# kubeadm join 10.0.50.10:6443 --token xflg2x.ehyvpkhdlfzm4n9h \
  # --discovery-token-ca-cert-hash sha256:3b6aacb2afe343b28bfa4c82def2fe68b0c065ed16fd35e5f2a0e2d42483c1d2
```

上面记录了完成的初始化输出的内容，根据输出的内容基本上可以看出手动初始化安装一个 Kubernetes 集群所需要的关键步骤。 其中有以下关键内容：

```text
* [certs]生成相关的各种证书
* [kubeconfig]生成相关的kubeconfig文件
* [kubelet-start] 生成kubelet的配置文件"/var/lib/kubelet/config.yaml"
* [control-plane]使用/etc/kubernetes/manifests目录中的yaml文件创建apiserver、controller-manager、scheduler的静态pod
* [bootstraptoken]生成token记录下来，后边使用kubeadm join往集群中添加节点时会用到
* 下面的命令是配置常规用户如何使用kubectl访问集群：mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config
* 最后给出了将节点加入集群的命令kubeadm join 10.0.50.10:6443 --token xflg2x.ehyvpkhdlfzm4n9h \
  --discovery-token-ca-cert-hash sha256:3b6aacb2afe343b28bfa4c82def2fe68b0c065ed16fd35e5f2a0e2d42483c1d2
```

查看一下集群状态，确认个组件都处于 healthy 状态:

```shell
kubectl get cs
# Warning: v1 ComponentStatus is deprecated in v1.19+
# NAME                 STATUS    MESSAGE                         ERROR
# scheduler            Healthy   ok
# controller-manager   Healthy   ok
# etcd-0               Healthy   {"health":"true","reason":""}
```

集群初始化如果遇到问题，可以使用 kubeadm reset 命令进行清理。

### 命令行自动补全

```shell
yum install -y bash-completion
kubectl completion bash >/etc/bash_completion.d/kubectl
# 重新进入shell即可
```

### 安装包管理器 helm 3

Helm 是 Kubernetes 的包管理器，后续流程也将使用 Helm 安装 Kubernetes 的常用组件。 这里先在 master 节点上安装 helm。

```shell
wget https://get.helm.sh/helm-v3.9.0-linux-amd64.tar.gz
tar -zxvf helm-v3.9.0-linux-amd64.tar.gz
mv linux-amd64/helm  /usr/local/bin/

helm list
```

### 部署 Pod Network 组件 Calico

选择 calico 作为 k8s 的 Pod 网络组件，下面使用 helm 在 k8s 集群中安装 calico。下载 tigera-operator 的 helm chart:

```shell
wget https://github.com/projectcalico/calico/releases/download/v3.26.1/tigera-operator-v3.26.1.tgz
```

查看这个 chart 的中可定制的配置:

```shell
helm show values tigera-operator-v3.23.1.tgz
```  

```yaml
imagePullSecrets: {}
  
installation:
  enabled: true
  kubernetesProvider: ""
  
apiServer:
  enabled: true
  
certs:
  node:
    key:
    cert:
    commonName:
  typha:
    key:
    cert:
    commonName:
    caBundle:
  
resources: {}
  
# Configuration for the tigera operator
tigeraOperator:
  image: tigera/operator
  version: v1.27.1
  registry: quay.io
calicoctl:
  image: docker.io/calico/ctl
  tag: v3.23.1
```

定制的 values.yaml 如下:

```text
# 可针对上面的配置进行定制,例如 calico 的镜像改成从私有库拉取。
# 这里只是个人本地环境测试 k8s 新版本，因此保留 value.yaml 为空即可
```

使用 helm 安装 calico：

```shell
helm install calico tigera-operator-v3.23.1.tgz -n kube-system  --create-namespace -f values.yaml
```

等待并确认所有 pod 处于 Running状态:

```shell
kubectl get pod -n kube-system | grep tigera-operator
# tigera-operator-5fb55776df-wxbph   1/1     Running   0             5m10s
kubectl get pods -n calico-system
# NAME                                       READY   STATUS    RESTARTS   AGE
# calico-kube-controllers-68884f975d-5d7p9   1/1     Running   0          5m24s
# calico-node-twbdh                          1/1     Running   0          5m24s
# calico-typha-7b4bdd99c5-ssdn2              1/1     Running   0          5m24s
```

查看一下 calico 向 k8s 中添加的 api 资源:

```shell
kubectl api-resources | grep calico
# bgpconfigurations                                                                 crd.projectcalico.org/v1               false        BGPConfiguration
# bgppeers                                                                          crd.projectcalico.org/v1               false        BGPPeer
# blockaffinities                                                                   crd.projectcalico.org/v1               false        BlockAffinity
# caliconodestatuses                                                                crd.projectcalico.org/v1               false        CalicoNodeStatus
# clusterinformations                                                               crd.projectcalico.org/v1               false        ClusterInformation
# felixconfigurations                                                               crd.projectcalico.org/v1               false        FelixConfiguration
# globalnetworkpolicies                                                             crd.projectcalico.org/v1               false        GlobalNetworkPolicy
# globalnetworksets                                                                 crd.projectcalico.org/v1               false        GlobalNetworkSet
# hostendpoints                                                                     crd.projectcalico.org/v1               false        HostEndpoint
# ipamblocks                                                                        crd.projectcalico.org/v1               false        IPAMBlock
# ipamconfigs                                                                       crd.projectcalico.org/v1               false        IPAMConfig
# ipamhandles                                                                       crd.projectcalico.org/v1               false        IPAMHandle
# ippools                                                                           crd.projectcalico.org/v1               false        IPPool
# ipreservations                                                                    crd.projectcalico.org/v1               false        IPReservation
# kubecontrollersconfigurations                                                     crd.projectcalico.org/v1               false        KubeControllersConfiguration
# networkpolicies                                                                   crd.projectcalico.org/v1               true         NetworkPolicy
# networksets                                                                       crd.projectcalico.org/v1               true         NetworkSet
# bgpconfigurations                 bgpconfig,bgpconfigs                            projectcalico.org/v3                   false        BGPConfiguration
# bgppeers                                                                          projectcalico.org/v3                   false        BGPPeer
# caliconodestatuses                caliconodestatus                                projectcalico.org/v3                   false        CalicoNodeStatus
# clusterinformations               clusterinfo                                     projectcalico.org/v3                   false        ClusterInformation
# felixconfigurations               felixconfig,felixconfigs                        projectcalico.org/v3                   false        FelixConfiguration
# globalnetworkpolicies             gnp,cgnp,calicoglobalnetworkpolicies            projectcalico.org/v3                   false        GlobalNetworkPolicy
# globalnetworksets                                                                 projectcalico.org/v3                   false        GlobalNetworkSet
# hostendpoints                     hep,heps                                        projectcalico.org/v3                   false        HostEndpoint
# ippools                                                                           projectcalico.org/v3                   false        IPPool
# ipreservations                                                                    projectcalico.org/v3                   false        IPReservation
# kubecontrollersconfigurations                                                     projectcalico.org/v3                   false        KubeControllersConfiguration
# networkpolicies                   cnp,caliconetworkpolicy,caliconetworkpolicies   projectcalico.org/v3                   true         NetworkPolicy
# networksets                       netsets                                         projectcalico.org/v3                   true         NetworkSet
# profiles                                                                          projectcalico.org/v3                   false        Profile
```

这些 api 资源是属于 calico 的，因此不建议使用 kubectl 来管理，推荐按照 calicoctl 来管理这些 api 资源。 将 calicoctl 安装为 kubectl 的插件:

```shell
cd /usr/local/bin
curl -o kubectl-calico -O -L  "https://github.com/projectcalico/calicoctl/releases/download/v3.21.5/calicoctl-linux-amd64" 
chmod +x kubectl-calico
kubectl-calico -h
```

### 验证 k8s DNS 是否可用

```shell
kubectl run curl --image=radial/busyboxplus:curl -it
# If you don't see a command prompt, try pressing enter.
# [ root@curl:/ ]$
```

进入后执行 nslookup kubernetes.default 确认解析正常:

```shell
nslookup kubernetes.default
# Server:    10.96.0.10
# Address 1: 10.96.0.10 kube-dns.kube-system.svc.cluster.local
# 
# Name:      kubernetes.default
# Address 1: 10.96.0.1 kubernetes.default.svc.cluster.local
```

### 向 Kubernetes 集群中添加 Node 节点

下面将 worker1, worker2 添加到 Kubernetes 集群中，分别在 worker1, worker2 上执行:

```shell
kubeadm join 10.0.50.10:6443 --token xflg2x.ehyvpkhdlfzm4n9h \
--discovery-token-ca-cert-hash sha256:3b6aacb2afe343b28bfa4c82def2fe68b0c065ed16fd35e5f2a0e2d42483c1d2
```

worker1 和 worker2 加入集群很是顺利，在 master 节点上执行命令查看集群中的节点：

```shell
kubectl get node
# NAME      STATUS   ROLES           AGE    VERSION
# master    Ready    control-plane   78m    v1.24.1
# worker1   Ready    <none>          95s    v1.24.1
# worker2   Ready    <none>          103s   v1.24.1
```

## Kubernetes 常用组件部署

### 使用 Helm 部署 ingress-nginx

为了便于将集群中的服务暴露到集群外部，需要使用 Ingress。接下来使用 Helm 将 ingress-nginx 部署到 Kubernetes上。 Nginx Ingress Controller 被部署在 Kubernetes 的边缘节点上。

这里将 master(10.0.50.10) 作为边缘节点，打上 Label：

```shell
kubectl label node master node-role.kubernetes.io/edge=
```

下载 ingress-nginx 的 helm chart:

```shell
wget https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-4.1.2/ingress-nginx-4.1.2.tgz
```

查看 ingress-nginx-4.1.2.tgz 这个 chart 的可定制配置:

```shell
helm show values ingress-nginx-4.1.2.tgz
```

对 values.yaml 配置定制如下:

```yaml
controller:
  ingressClassResource:
    name: nginx
    enabled: true
    default: true
    controllerValue: "k8s.io/ingress-nginx"
  admissionWebhooks:
    enabled: false
  replicaCount: 1
  image:
    # registry: k8s.gcr.io
    # image: ingress-nginx/controller
    # tag: "v1.2.0"
    registry: docker.io
    image: unreachableg/k8s.gcr.io_ingress-nginx_controller
    tag: "v1.2.0"
    digest: sha256:314435f9465a7b2973e3aa4f2edad7465cc7bcdc8304be5d146d70e4da136e51
  hostNetwork: true
  nodeSelector:
    node-role.kubernetes.io/edge: ''
  affinity:
    podAntiAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchExpressions:
            - key: app
              operator: In
              values:
              - nginx-ingress
            - key: component
              operator: In
              values:
              - controller
          topologyKey: kubernetes.io/hostname
  tolerations:
      - key: node-role.kubernetes.io/master
        operator: Exists
        effect: NoSchedule
      - key: node-role.kubernetes.io/master
        operator: Exists
        effect: PreferNoSchedule
```

nginx ingress controller 的副本数 replicaCount 为 1，将被调度到 master 这个边缘节点上。这里并没有指定 nginx ingress controller service 的 externalIPs，而是通过 hostNetwork: true 设置 nginx ingress controller 使用宿主机网络。 因为 k8s.gcr.io 被墙，这里替换成 unreachableg/k8s.gcr.io_ingress-nginx_controller 提前拉取一下镜像:

```shell
crictl pull unreachableg/k8s.gcr.io_ingress-nginx_controller:v1.2.0

helm install ingress-nginx ingress-nginx-4.1.2.tgz --create-namespace -n ingress-nginx -f values.yaml

kubectl get pod -n ingress-nginx
# NAME                                        READY   STATUS    RESTARTS   AGE
# ingress-nginx-controller-7f574989bc-xwbf4   1/1     Running   0          117s
```

测试访问 <http://10.0.50.10> 返回默认的 nginx 404 页，则部署完成。

### 使用 Helm 部署 dashboard

先部署 metrics-server：

```shell
wget https://github.com/kubernetes-sigs/metrics-server/releases/download/metrics-server-helm-chart-3.8.2/components.yaml
```

修改 components.yaml 中的 image 为 docker.io/unreachableg/k8s.gcr.io_metrics-server_metrics-server:v0.6.1。 修改 components.yaml 中容器的启动参数，加入 --kubelet-insecure-tls。

```shell
kubectl apply -f components.yaml
```

metrics-server 的 pod 正常启动后，等一段时间就可以使用 kubectl top 查看集群和 pod 的 metrics 信息:

```shell
kubectl top node
# NAME      CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
# master    871m         21%    3436Mi          44%
# worker1   239m         5%     1443Mi          18%
# worker2   282m         7%     1468Mi          19%
 
kubectl top pod -n kube-system
# NAME                               CPU(cores)   MEMORY(bytes)
# coredns-74586cf9b6-8klv8           6m           18Mi
# coredns-74586cf9b6-g2bzg           7m           16Mi
# etcd-master                        83m          86Mi
# kube-apiserver-master              228m         477Mi
# kube-controller-manager-master     62m          70Mi
# kube-proxy-56swk                   1m           19Mi
# kube-proxy-fzxdd                   8m           19Mi
# kube-proxy-tr4sp                   1m           18Mi
# kube-scheduler-master              11m          22Mi
# metrics-server-77cffb4988-26kp5    12m          20Mi
# tigera-operator-5fb55776df-dc6qz   12m          43Mi
```

接下来使用 helm 部署 k8s 的 dashboard，添加 chart repo:

```shell
helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/
helm repo update
```

查看 chart 的可定制配置:

```shell
helm show values kubernetes-dashboard/kubernetes-dashboard
```

对 values.yaml 定制配置如下:

```yaml
image:
  repository: kubernetesui/dashboard
  tag: v2.5.1
ingress:
  enabled: true
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
  hosts:
  - k8s.example.com
  tls:
    - secretName: example-com-tls-secret
      hosts:
      - k8s.example.com
metricsScraper:
  enabled: true
```

先创建存放 k8s.example.comssl 证书的 secret:

```shell
openssl genrsa -out myCA.key 2048
openssl req -x509 -new -nodes -key myCA.key -sha256 -days 1825 -out myCA.pem

kubectl create secret tls example-com-tls-secret \
  --cert=myCA.pem \
  --key=myCA.key \
  -n kube-system
```

使用 helm 部署 dashboard:

```shell
helm install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard \
-n kube-system \
-f values.yaml
```

确认上面的命令部署成功。创建管理员 sa:

```shell
kubectl create serviceaccount kube-dashboard-admin-sa -n kube-system
 
kubectl create clusterrolebinding kube-dashboard-admin-sa \
--clusterrole=cluster-admin --serviceaccount=kube-system:kube-dashboard-admin-sa
```

创建集群管理员登录 dashboard 所需 token:

```shell
kubectl create token kube-dashboard-admin-sa -n kube-system --duration=87600h

# eyJhbGciOiJSUzI1NiIsImtpZCI6IllCd2tjdnItclZuUXQ5ZjFjRUlJNjl3NFQzTGZPNGgwc1lCMzZfN0VfcjQifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxOTY5NDM3OTM0LCJpYXQiOjE2NTQwNzc5MzQsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJrdWJlLWRhc2hib2FyZC1hZG1pbi1zYSIsInVpZCI6IjBkMTcwYjg1LWE1YzItNDczOS1hZjI0LWQ2NjhlZmRjMzg2NCJ9fSwibmJmIjoxNjU0MDc3OTM0LCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06a3ViZS1kYXNoYm9hcmQtYWRtaW4tc2EifQ.KZU5kgbh1-wCSbP6mPoqZHhWWWGBicxBKasWU6Mq46Ujzx1pLAncE8vSP2ttwtzqPBYSfsqGLZMBjaIwlaBfHVNxZM_8MCeg39p7k9xfj7ljsWbh3o1ONUU8Pgikh1iFZFJuiBZkYeHrtwF4zIBd2KsosLqt3AUlF4FGciNTZESvhusB5LJRYkhT8i7oz78d3YRPVUd1x-J2H4z0fWzuziN7kPLh8X1tkBF7WxYdP2ppDdqMbCITvEr9RQoOemwe7jfqFQdURwgEV_iXrM6-o58rbwOFdciUFiI07dEg3r9uA552coK9YBAaRvpWxKgeCVcPhXu8gAHh-q8A_6PmFA
```

使用上面的 token 登录 k8s dashboard。
docker pull registry.k8s.io/kube-apiserver:v1.28.1
docker pull registry.k8s.io/kube-controller-manager:v1.28.1
docker pull registry.k8s.io/kube-scheduler:v1.28.1
docker pull registry.k8s.io/kube-proxy:v1.28.1
docker pull registry.k8s.io/pause:3.9
docker pull registry.k8s.io/etcd:3.5.9-0
docker pull registry.k8s.io/coredns/coredns:v1.10.1

docker save -o kube-apiserver_v1.28.1.tgz registry.k8s.io/kube-apiserver:v1.28.1
docker save -o kube-controller-manager_v1.28.1.tgz registry.k8s.io/kube-controller-manager:v1.28.1
docker save -o kube-scheduler_v1.28.1.tgz registry.k8s.io/kube-scheduler:v1.28.1
docker save -o kube-proxy_v1.28.1.tgz registry.k8s.io/kube-proxy:v1.28.1
docker save -o pause_3.9.tgz registry.k8s.io/pause:3.9
docker save -o etcd_3.5.9-0.tgz registry.k8s.io/etcd:3.5.9-0
docker save -o coredns_v1.10.1.tgz registry.k8s.io/coredns/coredns:v1.10.1

nerdctl --namespace k8s.io load -i coredns_v1.10.1.tgz
nerdctl --namespace k8s.io load -i etcd_3.5.9-0.tgz
nerdctl --namespace k8s.io load -i kube-apiserver_v1.28.1.tgz
nerdctl --namespace k8s.io load -i kube-controller-manager_v1.28.1.tgz
nerdctl --namespace k8s.io load -i kube-proxy_v1.28.1.tgz
nerdctl --namespace k8s.io load -i kube-scheduler_v1.28.1.tgz
nerdctl --namespace k8s.io load -i pause_3.9.tgz

export http_proxy="192.168.2.52:7890"
export https_proxy="192.168.2.52:7890"
export no_proxy="localhost, 127.0.0.1"