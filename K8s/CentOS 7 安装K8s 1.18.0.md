# CentOS 7 安装K8s 1.18.0

## 集群规划

| 角色   | IP               |
| ------ | ---------------- |
| master | d01(10.0.50.171) |
| node   | d02(10.0.50.172) |

## 系统初始化

```shell
# 关闭防火墙
systemctl disable firewalld
systemctl stop firewalld
# 关闭selinux
sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config
setenforce 0
# 关闭交换分区
sed -i '/ swap / s/^\(.*\)$/# \1/g' /etc/fstab
swapoff -a

```

## 每个节点添加hosts

```shell
cat >> /etc/hosts << EOF
10.0.50.171    k01
10.0.50.172    k02
EOF
```

## 每个节点安装Docker

安装docker [链接](https://blog.csdn.net/weixin_37714509/article/details/120189150)

## 每个节点安装kubeadm, kubelet, kubectl

```shell
# 添加阿里云的yum源
cat > /etc/yum.repos.d/kubernetes.repo << EOF
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF

# 版本更新太快，指定版本安装
yum install -y kubelet-1.24.0 kubeadm-1.24.0 kubectl-1.24.0 --exclude kubernetes

# 设置为开机自启动即可，由于没有生成配置文件，集群初始化后自动启动
systemctl enable kubelet
```

## Master 节点初始化

```shell
# 由于默认拉取镜像地址k8s.gcr.io国内无法访问，这里需要指定阿里云镜像仓库地址
kubeadm init \
  --apiserver-advertise-address=192.168.2.68 \
  --image-repository registry.aliyuncs.com/google_containers \
  --kubernetes-version v1.28.1 \
  --cri-socket=unix:///run/containerd/containerd.sock \
  --service-cidr=10.96.0.0/12 \
  --pod-network-cidr=10.244.0.0/16
```

## 根据提示进行操作, Master节点

```shell
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

wget https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

kubectl apply -f kube-flannel.yml
```

## Node节点

```shell
kubeadm join 10.0.50.171:6443 --token a573dm.82ej6kl69ghvwn6b \
    --discovery-token-ca-cert-hash sha256:fd69e6ead908b2faf117d3f8dd9b10b84d97112d0db893964f5a1fec1a7d3ae6
```

## 安装完成

```shell
kubectl get node
# NAME   STATUS   ROLES    AGE    VERSION
# k01    Ready    master   127m   v1.18.0
# k02    Ready    <none>   116m   v1.18.0
```

## 命令行自动补全

```shell
yum install -y bash-completion
kubectl completion bash >/etc/bash_completion.d/kubectl
# 重新进入shell即可
```

node-role.kubernetes.io/control-plane:NoSchedule
                    node.kubernetes.io/not-ready:NoSchedule