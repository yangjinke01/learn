# DNS

## DNS基础知识

### DNS的出现及演化

网络出现的早期是使用IP地址通讯的，那时就几台主机通讯。但是随着接入网络主机的增多，这种数字标识的地址非常不便于记忆，UNIX上就出现了建立一个叫做hosts的文件（Linux和Windows也继承保留了这个文件）。这个文件中记录着主机名称和IP地址的对应表。这样只要输入主机名称，系统就会去加载hosts文件并查找对应关系，找到对应的IP，就可以访问这个IP的主机了。
但是后来主机太多了，无法保证所有人都能拿到统一的最新的hosts文件，就出现了在文件服务器上集中存放hosts文件，以供下载使用。互联网规模进一步扩大，这种方式也不堪重负，而且把所有地址解析记录形成的文件都同步到所有的客户机似乎也不是一个好办法。这时DNS系统出现了，随着解析规模的继续扩大，DNS系统也在不断的演化，直到现今的多层架构体系。

### 什么是DNS

DNS（Domain Name
System，域名系统），因特网上作为域名和IP地址相互映射的一个分布式数据库，能够使用户更方便的访问互联网，而不用去记住能够被机器直接读取的IP数串。通过主机名，最终得到该主机名对应的IP地址的过程叫做域名解析（或主机名解析）。DNS可以使用TCP和UDP的53端口，基本使用UDP协议的53端口。DNS
的分布式数据库是以域名为索引的，每个域名实际上就是一棵很大的逆向树中路径，这棵逆向树称为域名空间（domain name space）。

### 域的分类

域是分层管理的。 第一层：根域

* 根域 ： .

第二层 ：顶级域 （tld top level domain）

* 按国家划分：.cn(中国)、.tw(台湾)、.hk(香港) 等。
* 按组织性质划分：.org、.net、.com、.edu、.gov 等。
* 反向域：arpa ，这是反向解析的特殊顶级域。

第三层及以下：

* 顶级域下来就是普通的域，公司或个人在互联网上注册的域名一般都是这些普通的域，如jd.com。

### DNS解析流程

![DNS解析流程](./images/dns.png)

这儿以我们访问 www.baidu.com 为例。

* 本地主机首先会查找本机 DNS 缓存，然後查询本地 hosts 文件是否有 www.baidu.com. 这个 FQDN 所对应的主机 IP 地址，若有，则直接使用；若没有，本机将向指定的 dns server 发起查询请求（这个
  DNS 服务器就是计算机里设置指向的 DNS）。
* DNS 服务器收到询问请求，首先查看自己是否有 www.baidu.com 的缓存，如果有就直接返回给客户端，没有就越级上访到根域"."，并询问根域。
* 根域只是记录了 .com 域的相关信息，所以将 .com 域的地址返回给 DNS 服务器。
* DNS 服务器根据根域返回的信息向 .com 域发起查询请求，由于 .com 域只记录了 baidu.com 的信息，所以将 baidu.com 域的地址返回给 DNS 服务器。
* DNS 服务器根据 .com 域返回的信息向 baidu.com 域发起查询请求，于是 baidu.com 域的 DNS 服务器就去查询本地的记录，找到了 www 主机对应 IP 地址，将该 IP 地址返回给 DNS 服务器。
* DNS 服务器将得到的 IP 地址返回给客户端，并缓存一份结果在自己机器中，方便下一次客户端再次访问该站点。
* 客户端得到回答的IP地址后缓存下来，并去访问 www.baidu.com，然后 www.baidu.com 就把页面内容发送给客户端，也就是百度页面。

注：

* 本机查找完缓存后如果没有结果，会先查找hosts文件，如果没有找到再把查询发送给DNS服务器，但这仅仅是默认情况，这个默认顺序是可以改变的。在/etc/nsswitch.conf中有一行" hosts: files dns"
  就是定义先查找hosts文件还是先提交给DNS服务器的，如果修改该行为"hosts: dns files"则先提交给DNS服务器，这种情况下hosts文件几乎就不怎么用的上了。
* 由于缓存是多层次缓存的，所以真正的查询可能并没有那么多步骤，上图的步骤是完全没有所需缓存的查询情况。假如某主机曾经向DNS服务器提交了www.baidu.com的查询，那么在DNS服务器上除了缓存了www.baidu.com的
  记录，还缓存了".com"和"baidu.com"的记录，如果再有主机向该DNS服务器提交ftp.baidu.com的查询，那么将跳过".“和”.com"的查询过程直接向baidu.com发出查询请求。
* DNS解析过程中存在两种查询类型：递归查询（从客户机至指定DNS服务器）、迭代查询（从DNS服务器至各个域）。

### DNS分类

* 主DNS服务器： 就是一台存储着原始资料的DNS服务器。
* 从DNS服务器： 使用自动更新方式从主DNS服务器同步数据的DNS服务器。也成辅助DNS服务器。
* 缓存服务器： 不负责本地解析，采用递归方式转发客户机查询请求，并返回结果给客户机的DNS服务器。同时缓存查询回来的结果，也叫递归服务器。
* 转发器： 这台DNS发现非本机负责的查询请求时，不再向根域发起请求，而是直接转发给指定的一台或者多台服务器。自身并不缓存查询结果。

### 资源记录

DNS服务器是如何根据主机名解析出 IP 地址，或从 IP 地址解析出主机名的呢？这儿我们就要用到资源记录（Resource Record），简称 RR。常用的记录类型有：A、AAAA、SOA、NS、PTR、CNAME、MX 等。

#### 资源记录的定义格式

```text
name   [TTL]   IN   RR_TYPE   value
```

##### SOA（Start Of Authority）： 起始授权记录，一个区域解析库有且只能有一个SOA记录，而且必须放在第一条

```text
name:当前域的名称，如 "baidu.com."，或"4.3.2.in-addr.arpa."。
value:由多部分组成。
(1) 当前域的名称（也可使用主DNS服务器名称）；
(2) 当前域管理员的邮箱地址，但地址中不能使用@符号，一般使用 "." 来代替；
(3) 主从服务协调属性的定义；
   第一个值是区域数据文件的序列编号serial，每次修改此区域数据文件都需要修改该编号值以便让slave dns服务器同步该区域数据文件。
   第二个值是刷新refresh时间间隔，表示slave dns服务器找master dns服务器更新区域数据文件的时间间隔。
   第三个值是重试retry时间间隔，表示slave dns服务器找master dns服务器更新区域数据文件时，如果联系不上master，则等待多久再重试联系，该值一般比refresh时间短，否则该值表示的重试就失去了意义。
   第四个值是过期expire时间值，表示slave dns服务器上的区域数据文件多久过期。
   第五个值是negative answer ttl，表示客户端找dns服务器解析时，否定答案的缓存时间长度。
   这几个值可以分行写，也可以直接写在同一行中使用空格分开。

例如：
 test.com. IN  SOA  test.com.  admin.test.com.  (
          2018110601 ;serial
          2H         ;refresh  2 hours
          10M        ;retry  10 min
          1W         ;expire   1 week
          1D         ;negative answer ttl  1 day
 )
```

##### NS（Name Server）：存储的是该域内的 DNS 服务器相关信息。即 NS 记录标识了哪台服务器是 DNS 服务器

```text
name:当前域的名称。
value:当前域的某 DNS 服务器的名称，如 "ns.test.com."。

例如
test.com  IN  NS  ns1.test.com
test.com  IN  NS  ns2.test.com

注：一个域内可以有多个 ns 记录，即可以存在多台 DNS 服务器。
```

##### A（Address）：存储的是域内主机名所对应的ip地址

```text
name:某 FQDN，如 "www.test.com."。
value:某 IPv4 地址。

例如：
www.test.com.  IN  A  192.168.100.200

注：AAAA记录格式和A记录格式相似，但 value 是某 IPv6 地址。
```

##### PTR（Pointer）：和A记录相反，存储的是 ip 地址对应的主机名，该记录只存在于反向解析的区域数据文件中(并非一定)

```text
name:IP 地址，有特定格式，且加上特定后缀，如："1.2.3.4" 的记录应该写为 "4.3.2.in-addr.arpa"。 value:某 FQDN

例如： 4.3.2.in-addr.arpa IN PTR  www.test.com.
```

##### CNAME（Canonical Name）：表示规范名的意思，其所代表的记录常称为别名记录。之所以如此称呼，就是因为为规范名起了一个别名。什么是规范名？可以简单认为是 FQDN

```text
name: FQDN 格式的别名 value: FQDN 格式的初始名

例如： web.test.com. IN CNAME  www.test.com
```

##### MX（Mail Exchanger）：邮件交换器

```text
name:当前的域名 value:当前域内某邮件交换器的主机名

例如
test.com. IN MX 10 mx1.test.com.
test.com. IN MX 20 mx2.test.com.

注：MX记录可以有多个，但每一个记录的value之前应该有一个数字表示优先级 。优先级：0-99，数字越小优先级越高。
```

## DNS安装配置

### 安装DNS

Bind 是一款开放源码的 DNS 服务器软件，Bind由美国加州大学 Berkeley 分校开发和维护的，全名为 Berkeley Internet Name Domain 它是目前世界上使用最为广泛的 DNS。
**实验环境：CentOS 7**

```shell
yum install -y bind
```

### 配置文件解析

#### 配置文件列表

```text
主配置文件：/etc/named.conf
    主配置文件包含进来的其他文件：
        /etc/named.iscdlv.key
            /etc/named.rfc1912.zones
            /etc/named.root.key
    解析库文件
        /var/named/目录下：一般名字为：ZONE_NAME.zone
```

#### 主配置文件：named.conf

```text
options {          ---- 全局配置段
 ...
};

logging {   ---- 日志配置段
 ...
};

zone "." IN {    ---- 区域配置段，可定义在主配置文件，也可定义在"/etc/named.rfc1912.zones"文件中
        type hint;
        file "named.ca";
};

include "/etc/named.rfc1912.zones";
include "/etc/named.root.key";

注：每个配置语句必须以分号结尾。
```

#### 配置正向解析

这里以解析 jack.com 域为例。

* 修改主配置文件：

```shell
vim /etc/named.conf
#    listen-on port 53 { 10.0.50.16; };    # DNS服务器地址
```

* 修改 named.rfc1912.zones 文件：

```text
~]# vim /etc/named.rfc1912.zones

zone "jack.com" IN {        ----> 域定义格式
    type master;            ---> [hint|master|slave|forward] 根,主,从,转发
    file "jack.com.zone";   ----> 自定义解析域文件名称
};
#添加到最后即可
```

* 创建"jack.com.zone"文件：

```text
~]# cd /var/named/
~]# vim jack.com.zone
$TTL 1D
@       IN SOA  @ yangjinke80.gmail.com. (
                                        20220720        ; serial
                                        1D              ; refresh
                                        1H              ; retry
                                        1W              ; expire
                                        3H )            ; minimum
        NS      ns1
ns1     A       10.0.50.16
*       A       10.0.50.16

~]# chown root:named jack.com.zone
~]# chmod 640 jack.com.zone

注：
(1) "$"符号：定义宏。最常见的是"$TTL"、"$ORIGIN"。
(2) FQDN自动补齐：在区域数据文件中，没有使用点号"."结尾的，在实际使用的时候都会自动补上域名，使其变为 FQDN。
    例如：上面文件中的 ns1,会自动补全为 ns1.chuan.com.
(3) 若上一条记录与下一条记录主机名相同，则下一条可以省略，默认为上一条的主机名。如上文件中 www下为空，默认为下一条主机名也为www。
(4) "@" 默认代表代表域名。
```

* 检查配置文件

```text
~]# named-checkconf        #默认检查 named.conf 和 named.rfc1912.zones 文件
~]# named-checkzone jack.com /var/named/jack.com.zone
zone jack.com/IN: loaded serial 20220720
OK
```

* 启动服务

```shell
systemctl start named
```

* 修改本机DNS服务器指向

现在我们知道，一台主机访问另一台主机，需要通过FQDN 解析出对应的 IP，首先会先查找本地缓存以及 hosts文件，若没有就会向本机指向的 DNS 服务器发起递归请求。所有此处我们需要将本机指向的 DNS 服务器指向我们自己搭建的 DNS 服务器，否则将无法被解析。

```shell
~]# vim /etc/resolv.conf
nameserver 10.0.50.16
```

* 测试DNS解析

常用的测试DNS解析的命令有 nslookup、host、dig。这儿以 dig 进行测试。若没有 dig 工具，可装上 bind-utils 包。
dig：

```text
用法：
    dig [-t RR_TYPE] name [@SERVER] [query options]
正向解析：
    dig -t A name [@SERVER]
反向解析：
    dig -x IP
```

```text
dig -t A photo.jack.com @10.0.50.16

...
;; ANSWER SECTION:
photo.jack.com.        86400    IN    A    10.0.50.16
...
```

### DNS访问控制

访问控制是指仅对定义的网络进行解析。访问控制是通过 acl 函数来实现的，acl 把一个或多个地址归并为一个集合，并通过一个统一的名称调用。需要注意的是：acl 只能先定义，后使用。因此，其一般在 named.conf 文件的 options 字段的前面定义。

acl 的格式：

```text
acl    acl_name {
    ip；具体的ip地址
    net/prelen；表示一个网段
};

例如：
acl mynet {
    172.168.179.110
    172.168.179.0/24
};

allow-query     { localhost; };
allow-query     { mynet; };
...
```

bind 内置的 acl：

```text
none    ：没有一个主机
any        ：任意主机
localhost    ：本地主机
localnets：本机的IP同掩码 
```

<https://blog.csdn.net/rightlzc/article/details/83756810>
