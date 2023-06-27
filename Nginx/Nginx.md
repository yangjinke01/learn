# Nginx

## Compile Install

[download](https://nginx.org/en/download.html)

```shell
tar -xf nginx-1.22.0.tar.gz
cd nginx-1.22.0

yum install -y gcc pcre pcre-devel zlib zlib-devel make
# apt-get install build-essential libpcre3 libpcre3-dev zlib1g zlib1g-dev libssl-dev libgd-dev libxml2 libxml2-dev uuid-dev

./configure --help

./configure --prefix=/opt/nginx

make && make install

ln -s /opt/nginx/sbin/nginx /usr/local/sbin/nginx

nginx -h
```

## Service Configuration

```shell
cat > /usr/lib/systemd/system/nginx.service << EOF
[Unit]
Description=The nginx HTTP and reverse proxy server
After=network-online.target remote-fs.target nss-lookup.target

[Service]
Type=forking
PIDFile=/opt/nginx/logs/nginx.pid
ExecStartPre=/usr/bin/rm -f /opt/nginx/logs/nginx.pid
ExecStartPre=/opt/nginx/sbin/nginx -t -c /opt/nginx/conf/nginx.conf
ExecStart=/opt/nginx/sbin/nginx -c /opt/nginx/conf/nginx.conf
ExecReload=/opt/nginx/sbin/nginx -s reload
ExecStop=/opt/nginx/sbin/nginx -s stop
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl start nginx
systemctl status nginx
```

## Config file

### Default configuration

```nginx configuration
# 进程数，建议等于cpu核心数
worker_processes  1;

events {
    # 每个进程处理的连接数
    worker_connections  1024;
}

http {
    # 包括其它的配置文件
    include       mime.types; # 根据文件后缀确定文件类型，响应头

    # 如果没在mime.types里面，取默认值
    default_type  application/octet-stream;

    # 零拷贝，不用加载数据到nginx，nginx发信号给内核直接发送数据
    sendfile        on;

    # 保持连接的时间
    keepalive_timeout  65;

    # 虚拟主机，相互隔离
    server {
        # 根据端口和域名区分主机
        listen       80;
        # 域名或主机名
        server_name  localhost;
        # URI
        location / {
            # 相对于nginx安装目录
            root   html;
            # 默认主页
            index  index.html index.htm;
        }
        # 服务端错误重定向到该URI
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
    }
}
```

### Virtual Server Name

需要配置DNS

```nginx configuration
worker_processes  1;
events {
    worker_connections  1024;
}
http {
    include       mime.types;
    default_type  application/octet-stream;
    sendfile        on;
    keepalive_timeout  65;
    server {
        listen      80;
        server_name  localhost;
        location / {
            root   html;
            index  index.html index.htm;
        }
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
    }
    server {
        listen      80;
        server_name  html.jack.com;
        location / {
            root   /data/www/html;
            index  index.html;
        }
    }
    server {
        listen      80;
        server_name photo.jack.com photos.jack.com;
        location / {
            root    /data/www/photo;
            index   index.html;
        }
    }
    server {
        listen      80;
        server_name *.jack.com;
        location / {
            root    /data/www/video;
            index   index.html;
        }
    }
}
```

### Reverse Proxy

```nginx configuration
events {
    worker_connections  1024;
}
http {
    server {
        location / {
            proxy_pass http://www.baidu.com;
        }
    }
}
```

### Load Balance

```nginx configuration
events {
    worker_connections  1024;
}

http {
    upstream workers {
        # = 两边不能有空格
        server worker1.jack.com weight=3;
        # [down(不参与调度) | backup(其它主机都挂了才调度)]
        server worker2.jack.com weight=1;
    }
    
    server {
        location / {
            proxy_pass http://workers;
        }
    }
}
```

### 动静分离、URL重写

<https://blog.csdn.net/zxd1435513775/article/details/102508549>

