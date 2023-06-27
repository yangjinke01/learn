# 脚本语法

### 脚本应该放的位置

```shell
# 去掉后缀，因为脚本的第一行是shell bell，可以执行
/usr/local/bin/
```



```shell
# 没有该变量，则设置为2
KUBE_VERBOSE="${KUBE_VERBOSE:-2}"

# 匹配符合规则的行，并打印指定字段
awk -F '=' '/PRETTY_NAME/ {print $2}' /etc/os-release

# 取字符串的某个范围
OS=$(uname -a)
echo ${OS:0:13}

# 模块
lsmod
modprobe br_netfilter

# Load settings from all system configuration files
sysctl --system
```

## Install iftop

```shell
yum install epel-release
yum install  iftop
```

## read

```shell
while true; do
    read -p "Do you wish to install this program? " yn
    case $yn in
        [Yy]* ) make install; break;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done
```

### if
```shell
# 中括号用于test, man test
if command -v $command
```

### expr

```shell
# 特殊符号需要加转义符
expr 3 \* 2
```

### command

```shell
# 判断某个命令是否存在
command -v htop
```

### for

```shell
for file in logfiles/*.log;do ls $file;done
for num in {1..10};do echo $num;done
```

### data stream

```shell
2> stderr
1> stdout 1 可以省略
&> stderr and stdout

一条命令中可以同时用多个重定向
<command> 2> error.txt 1> success.txt
```

