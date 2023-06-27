# CentOS 安装 Docker

## 安装

```shell
# yum-config-manager命令所在安装包
yum install -y yum-utils

# 添加阿里云的docker yum源
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

# 安装最新版docker
yum install -y docker-ce-18.09.1

# 安装历史版本
# yum list docker-ce --showduplicates
# 从第一个冒号（:）一直到第一个连字符，并用连字符（-）分隔。例如：docker-ce-18.09.1
# yum install docker-ce-<VERSION_STRING>
```

## 启动

```shell
# 开机自启动
systemctl enable docker
# 启动和查看状态
systemctl start docker
systemctl status docker
```

## 卸载

```shell
# 停止
systemctl disable docker
systemctl stop docker

yum remove docker-ce
# 删除镜像、容器、配置文件等内容
rm -rf /var/lib/docker
```

## 强制删除镜像

```shell
# [root@gw04 ~]# docker images
# REPOSITORY             TAG       IMAGE ID       CREATED         SIZE
# gitlab/gitlab-ee       latest    cea68caaec5b   2 months ago    2.78GB

rm /var/lib/docker/image/devicemapper/imagedb/content/sha256/cea68caaec5b*
```

## 配置国内加速

```shell
mkdir /etc/docker
cat > /etc/docker/daemon.json << EOF
{
  "registry-mirrors": [
    "http://hub-mirror.c.163.com",
    "https://docker.mirrors.ustc.edu.cn"
  ]
}
EOF

systemctl daemon-reload
systemctl restart docker
```
