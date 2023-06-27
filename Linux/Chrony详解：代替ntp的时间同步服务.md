# Chrony详解：代替ntp的时间同步服务

## 1.chrony简介

------

Chrony是一个开源的自由软件，它能保持系统时钟与时钟服务器（NTP）同步，让时间保持精确。

它由两个程序组成：chronyd和chronyc。

chronyd：是守护进程，主要用于调整内核中运行的系统时间和时间服务器同步。它确定计算机增减时间的比率，并对此进行调整补偿。

chronyc：提供一个用户界面，用于监控性能并进行多样化的配置。它可以在chronyd实例控制的计算机上工作，也可以在一台不同的远程计算机上工作。

chrony是[CentOS](https://www.linuxidc.com/topicnews.aspx?tid=14)7.x上自带的时间同步软件，chrony既可以做服务端，又可以做客户端，ntp相对来说更简单一些。

**chrony相对于ntp的优势：**

- 更快的同步只需要数分钟而非数小时时间，从而最大程度减少了时间和频率误差，这对于并非全天 24 小时运行的台式计算机或系统而言非常有用。
- 能够更好地响应时钟频率的快速变化，这对于具备不稳定时钟的虚拟机或导致时钟频率发生变化的节能技术而言非常有用。
- 在初始同步后，它不会停止时钟，以防对需要系统时间保持单调的应用程序造成影响。
- 在应对临时非对称延迟时（例如，在大规模下载造成链接饱和时）提供了更好的稳定性。
- 无需对服务器进行定期轮询，因此具备间歇性网络连接的系统仍然可以快速同步时钟。

## 2.CentOS7安装chrony

------

1.CentOS7系统默认已经安装，如未安装，请执行以下命令安装：

```shell
yum install chrony -y
```

2.启动并加入开机自启动

```shell
systemctl enable --now chronyd
```

3.Firewalld设置

```shell
firewall-cmd --add-service=ntp --permanent
firewall-cmd --reload
```

因NTP使用123/UDP端口协议，所以允许NTP服务即可。

4.设置时区

```shell
# 查看当前系统时区：
timedatectl
#       Local time: Fri 2018-2-29 13:31:04 CST
#   Universal time: Fri 2018-2-29 05:31:04 UTC
#         RTC time: Fri 2018-2-29 08:17:20
#        Time zone: Asia/Shanghai (CST, +0800)
#      NTP enabled: yes
# NTP synchronized: yes
#  RTC in local TZ: no
#       DST active: n/a

# 如果你当前的时区不正确，请按照以下操作设置。

# 查看所有可用的时区：

timedatectl list-timezones

# 筛选式查看在亚洲S开的上海可用时区：

timedatectl list-timezones |  grep  -E "Asia/S.*"

# Asia/Sakhalin
# Asia/Samarkand
# Asia/Seoul
# Asia/Shanghai
# Asia/Singapore
# Asia/Srednekolymsk

# 设置当前系统为Asia/Shanghai上海时区：

timedatectl set-timezone Asia/Shanghai

# 设置完时区后，强制同步下系统时钟：

chronyc -a makestep
# 200 OK

# 查看时间同步源：
chronyc sources -v

# 查看时间同步源状态：
chronyc sourcestats -v

# 设置硬件时间，硬件时间默认为UTC：
timedatectl set-local-rtc 1

# 启用NTP时间同步：
timedatectl set-ntp yes

# 校准时间服务器：
chronyc tracking
```

## 配置

### 查看配置文件目录

```shell
rpm -ql chrony
# /etc/NetworkManager/dispatcher.d/20-chrony
# /etc/chrony.conf
# /etc/chrony.keys
# /etc/dhcp/dhclient.d/chrony.sh
# /etc/logrotate.d/chrony
# /etc/rc.d/init.d/chronyd
# /etc/sysconfig/chronyd
# /usr/bin/chronyc
# /usr/sbin/chronyd
# /usr/share/doc/chrony-2.1.1
# /usr/share/doc/chrony-2.1.1/COPYING
# /usr/share/doc/chrony-2.1.1/FAQ
# /usr/share/doc/chrony-2.1.1/NEWS
# /usr/share/doc/chrony-2.1.1/README
# /usr/share/doc/chrony-2.1.1/chrony.txt.gz
# /usr/share/info/chrony.info.gz
# /usr/share/man/man1/chronyc.1.gz
# /usr/share/man/man5/chrony.conf.5.gz
# /usr/share/man/man8/chronyd.8.gz
# /var/lib/chrony
# /var/lib/chrony/drift
# /var/lib/chrony/rtc
# /var/log/chrony
```

### 查看chrony配置

```shell
# cat /etc/chrony.conf
# Use public servers from the pool.ntp.org project.
# Please consider joining the pool (http://www.pool.ntp.org/join.html).
server 0.rhel.pool.ntp.org iburst
server 1.rhel.pool.ntp.org iburst
server 2.rhel.pool.ntp.org iburst
server 3.rhel.pool.ntp.org iburst

# Ignore stratum in source selection.
stratumweight 0

# Record the rate at which the system clock gains/losses time.
driftfile /var/lib/chrony/drift

# In first three updates step the system clock instead of slew
# if the adjustment is larger than 10 seconds.
makestep 10 3

# Enable kernel synchronization of the real-time clock (RTC).
rtcsync

# Allow NTP client access from local network.
#allow 192.168/16

# Serve time even if not synchronized to any NTP server.
#local stratum 10

# Specify file containing keys for NTP and command authentication.
keyfile /etc/chrony.keys

# Specify key number for command authentication.
commandkey 1

# Generate new command key on start if missing.
generatecommandkey

# Disable logging of client accesses.
noclientlog

# Send message to syslog when clock adjustment is larger than 0.5 seconds.
logchange 0.5

# Specify directory for log files.
logdir /var/log/chrony

# Select which information is logged.
#log measurements statistics tracking
```

### chrony.conf配置参数说明

| 参数              | 参数说明                                                     |
| ----------------- | ------------------------------------------------------------ |
| **server**        | 该参数可以多次用于添加时钟服务器，必须以"server "格式使用。一般而言，你想添加多少服务器，就可以添加多少服务器 |
| **stratumweight** | stratumweight指令设置当chronyd从可用源中选择同步源时，每个层应该添加多少距离到同步距离。默认情况下，CentOS中设置为0，让chronyd在选择源时忽略源的层级 |
| **driftfile**     | chronyd程序的主要行为之一，就是根据实际时间计算出计算机增减时间的比率，将它记录到一个文件中是最合理的，它会在重启后为系统时钟作出补偿，甚至可能的话，会从时钟服务器获得较好的估值 |
| **makestep**      | 通常，chronyd将根据需求通过减慢或加速时钟，使得系统逐步纠正所有时间偏差。在某些特定情况下，系统时钟可能会漂移过快，导致该调整过程消耗很长的时间来纠正系统时钟。该指令强制chronyd在调整期大于某个阀值时步进调整系统时钟，但只有在因为chronyd启动时间超过指定限制（可使用负值来禁用限制），没有更多时钟更新时才生效。`maketep 1 3` 如果它的偏移大于一秒，则时钟将在前三次更新中校准。通常，建议仅在前几次更新中允许该步骤，但在某些情况下（例如，没有RTC或虚拟机的计算机可以在不正确的时间内暂停和恢复）可能需要允许任何步骤时钟更新。上面的例子可以改为 `maketep 1 -1` |
| **rtcsync**       | rtcsync指令将启用一个内核模式，在该模式中，系统时间每11分钟会拷贝到实时时钟（RTC） |
| **allow/deny**    | 这里你可以指定一台主机、子网，或者网络以允许或拒绝NTP连接到扮演时钟服务器的机器 |
| **keyfile**       | 指定包含NTP验证密钥的文件                                    |
| **logchange**     | 该指令设置调整系统时钟的阈值，超过该阈值将生成系统日志消息。通过NTP数据包，参考时钟或通过chronyc的settime命令输入的时间戳检测到时钟错误。默认情况下，阈值为1秒。使用的一个例子是：`logchange 0.1`如果开始补偿超过0.1秒的系统时钟错误，将导致生成系统日志消息。 |
| **logdir**        | 指定日志文件的目录。                                         |

### 配置chrony.conf

```text
server pool.ntp.org iburst
stratumweight 0
driftfile /var/lib/chrony/drift
makestep 1 3
rtcsync
keyfile /etc/chrony.keys
commandkey 1
generatecommandkey
noclientlog
logdir /var/log/chrony
```

### 设置时区，启动chrony

```shell
date
# Sun Dec 30 12:01:14 CST 2018

cat /etc/sysconfig/clock 
# ZONE="Asia/Shanghai"
# UTC=false
# ARC=false

/etc/init.d/chronyd start

chronyc sourcestats
# 210 Number of sources = 1
# Name/IP Address            NP  NR  Span  Frequency  Freq Skew  Offset  Std Dev
# ==============================================================================
# ntp1.ams1.nl.leaseweb.net   4   3     7    -29.989   4320.438  -1166us   649us

chronyc sources -v
# 210 Number of sources = 1

#   .-- Source mode  '^' = server, '=' = peer, '#' = local clock.
#  / .- Source state '*' = current synced, '+' = combined , '-' = not combined,
# | /   '?' = unreachable, 'x' = time may be in error, '~' = time too variable.
# ||                                                 .- xxxx [ yyyy ] +/- zzzz
# ||      Reachability register (octal) -.           |  xxxx = adjusted offset,
# ||      Log2(Polling interval) --.      |          |  yyyy = measured offset,
# ||                                \     |          |  zzzz = estimated error.
# ||                                 |    |           \
# MS Name/IP address         Stratum Poll Reach LastRx Last sample
# ===============================================================================
# ^* ntp1.ams1.nl.leaseweb.net     2   6    17    49   -472us[-1469us] +/-  277ms
```

### chronyc命令参数说明

| **参数**        | **参数说明**                     |
| --------------- | -------------------------------- |
| **sources**     | 查看时间同步源                   |
| **sourcestats** | 查看时间同步源状态               |
| **accheck**     | 检查NTP访问是否对特定主机可用    |
| **activity**    | 该命令会显示有多少NTP源在线/离线 |
| **add server**  | 手动添加一台新的NTP服务器。      |
| **clients**     | 在客户端报告已访问到服务器       |
| **delete**      | 手动移除NTP服务器或对等服务器    |
| **settime**     | 手动设置守护进程时间             |
| **tracking**    | 显示系统时间信息                 |

**注意：**

1.需要注意的是，配置完/etc/chrony.conf后，需重启chrony服务，否则可能会不生效。

2.chrony与ntp都是时间同步软件，两者最好不要同时开启。

3.若内网时间同步服务器有问题，直接把同步程序关掉，`ntpdate host_ip` 手动指向要同步时间的服务器即可。
