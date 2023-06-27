# 二进制安装k8s

## 注意事项

* 不要使用中文服务器

* 不要用克隆的虚拟机

* 更换本文档IP请全局替换

## IP 规划

| 角色          | IP         |
| ------------- | ---------- |
| k8s-master-01 | 10.0.50.11 |
| k8s-master-02 | 10.0.50.12 |
| k8s-master-01 | 10.0.50.13 |
| k8s-node-01   | 10.0.50.14 |
| k8s-node-02   | 10.0.50.15 |
| k8s-master-lb | 10.0.50.16 |

## 环境准备
### Ansible 资产文件
```ini
# Ansible Inventory k8s_hosts.inventory

[master]
k8s-master-01 ansible_host=k8s-master-01 ansible_user=root
k8s-master-02 ansible_host=k8s-master-02 ansible_user=root
k8s-master-03 ansible_host=k8s-master-03 ansible_user=root

[node]
k8s-node-01 ansible_host=k8s-node-01 ansible_user=root
k8s-node-02 ansible_host=k8s-node-02 ansible_user=root
```
### 修改主机名
```shell
ansible -i k8s_hosts.inventory k8s-master-01 -m hostname -a name=k8s-master-01
ansible -i k8s_hosts.inventory k8s-master-02 -m hostname -a name=k8s-master-02
ansible -i k8s_hosts.inventory k8s-master-03 -m hostname -a name=k8s-master-03
ansible -i k8s_hosts.inventory k8s-node-01 -m hostname -a name=k8s-node-01
ansible -i k8s_hosts.inventory k8s-node-02 -m hostname -a name=k8s-node-02
```

### docker yum 源 docker-ce-stable.repo

```ini
# docker-ce-stable.repo
[docker-ce-stable]
name=Docker CE Stable - $basearch
baseurl=https://mirrors.aliyun.com/docker-ce/linux/centos/$releasever/$basearch/stable
enabled=1
gpgcheck=1
gpgkey=https://mirrors.aliyun.com/docker-ce/linux/centos/gpg
```

```shell
ansible -i k8s_hosts.inventory all -m copy -a 'src=./docker-ce-stable.repo dest=/etc/yum.repos.d/'
```

### 安装Docker,禁用swap

install_docker.yml

```yaml
---
- hosts: "all"
  tasks:
    - name: Disable swap
      shell: |
        swapoff -a
        sed -i '/ swap / s/^\(.*\)$/# \1/g' /etc/fstab
    - name: Add Docker YUM repository
      shell: |
        yum install -y yum-utils
        yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

    - name: Install Docker
      yum:
        name: docker-ce
        state: present
```

```shell
ansible-playbook -i k8s_hosts.inventory install_rpm.yml
```

### 关闭防火墙和selinux

```shell
ansible -i k8s_hosts.inventory all -m command -a "systemctl disable --now firewalld"
ansible -i k8s_hosts.inventory all -m command -a "setenforce 0"
ansible -i k8s_hosts.inventory all -m command -a "sed -i s/^SELINUX=enforcing/SELINUX=disabled/ /etc/selinux/config"
```

### 时间同步

```shell
ansible -i k8s_hosts.inventory all -m yum -a "name=chrony state=present"
ansible -i k8s_hosts.inventory all -m command -a "systemctl enable --now chronyd"
```

### 升级内核

```shell
yum update -y --exclude=kernel*
ansible -v -i k8s_hosts.inventory all -m command -a "yum localinstall -y /root/kernel-lt-5.4.227-1.el7.elrepo.x86_64.rpm"
ansible -v -i k8s_hosts.inventory all -m command -a "yum localinstall -y /root/kernel-lt-devel-5.4.227-1.el7.elrepo.x86_64.rpm"

ansible -v -i k8s_hosts.inventory all -m command -a "grub2-set-default 0"
ansible -v -i k8s_hosts.inventory all -m command -a "grub2-mkconfig -o /boot/grub2/grub.cfg"
ansible -v -i k8s_hosts.inventory all -m command -a "reboot"
```

### IPVS

```shell
ansible -v -i k8s_hosts.inventory all -m yum -a "name=ipset,ipvsadm state=present"

# ipvs.modules 文件内容

#!/bin/bash
modprobe -- ip_vs
modprobe -- ip_vs_rr
modprobe -- ip_vs_wrr
modprobe -- ip_vs_sh
modprobe -- nf_conntrack

ansible -v -i k8s_hosts.inventory all -m copy -a "src=./ipvs.modules dest=/etc/sysconfig/modules/ mode=755"
ansible -v -i k8s_hosts.inventory all -m command -a "bash /etc/sysconfig/modules/ipvs.modules"
```

### 核心转发

```shell
ansible -v -i k8s_hosts.inventory all -m copy -a "content='net.ipv4.ip_forward = 1
user.max_user_namespaces = 28633
vm.swappiness = 0
' dest=/etc/sysctl.d/99-kubernetes-cri.conf mode=755"

ansible -v -i k8s_hosts.inventory all -m command -a "sysctl -p /etc/sysctl.d/99-kubernetes-cri.conf"
```

