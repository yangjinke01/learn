# CentOS7 升级内核

## 联网升级内核

### 1. 查看内核版本

```bash
uname -r 　　 
```

### 2. 导入ELRepo软件仓库的公共秘钥

```bash
rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org
```

### 3. 安装ELRepo软件仓库的yum源

```bash
rpm -Uvh http://www.elrepo.org/elrepo-release-7.0-3.el7.elrepo.noarch.rpm
```

### 4. 启用 elrepo 软件源并下载安装最新长期支持版内核

```bash
yum --enablerepo=elrepo-kernel install kernel-lt -y
```

### 5. 查看系统可用内核，并设置内核启动顺序

```bash
sudo awk -F\' '$1=="menuentry " {print i++ " : " $2}' /etc/grub2.cfg
```

### 6. 生成 grub 配置文件

```bash
# 机器上存在多个内核，我们要使用最新版本，可以通过 grub2-set-default 0 命令生成 grub 配置文件
grub2-set-default 0 　　#初始化页面的第一个内核将作为默认内核
grub2-mkconfig -o /boot/grub2/grub.cfg　　#重新创建内核配置
```

### 7. 重启系统并验证

```bash
reboot
uname -r
```

### 8. 删除旧内核

```bash
yum -y remove kernel kernel-tools
```

## 离线升级内核

### 1. 到官网下载内核 rpm 包

官方地址：<https://elrepo.org/linux/kernel/el7/x86_64/RPMS/>

下载最新版的内核并上传至服务器

### 2. 执行安装

```bash
rpm -ivh kernel-lt-5.4.95-1.el7.elrepo.x86_64.rpm kernel-lt-devel-5.4.95-1.el7.elrepo.x86_64.rpm
```

### 3. 查看系统可用内核，并设置内核启动顺序

```bash
sudo awk -F\' '$1=="menuentry " {print i++ " : " $2}' /etc/grub2.cfg
```

### 4. 生成 grub 配置文件

```bash
grub2-set-default 0 && grub2-mkconfig -o /boot/grub2/grub.cfg
```

### 5. 重启系统，并验证

```bash
reboot
uname -r
```

### 6. 删除旧内核

```bash
yum -y remove kernel kernel-tools
```
