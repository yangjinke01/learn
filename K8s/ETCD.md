# ETCD

## Prepare cert

```shell
cat > config.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "876000h"
    },
    "profiles": {
      "kubernetes": {
        "usages": ["signing", "key encipherment", "server auth", "client auth"],
        "expiry": "876000h"
      },
      "client": {
        "usages": ["signing", "key encipherment", "client auth"],
        "expiry": "876000h"
      },
      "peer": {
        "usages": ["signing", "key encipherment", "server auth", "client auth"],
        "expiry": "876000h"
      }
    }
  }
}
EOF


cat > ca-csr.json <<EOF
{
  "CN": "root-ca",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "k8s",
      "OU": "dyrnq"
    }
  ],
  "ca": {
    "expiry": "876000h"
 }
}
EOF

cat > server-csr.json <<EOF
{
  "CN": "etcd",
  "hosts": [
    "ubuntu",
    "192.168.33.11",
    "192.168.33.12",
    "192.168.33.13",
    "localhost",
    "127.0.0.1"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "k8s",
      "OU": "dyrnq"
    }
  ]
}
EOF
cat > client-csr.json <<EOF
{
  "CN": "etcd",
  "key": {
    "algo": "rsa",
    "size": 4096
  },
  "names": [
    {
      "O": "etcd"
    }
  ]
}
EOF
cat > peer-csr.json <<EOF
{
  "CN": "peer",
  "hosts": [
    "ubuntu",
    "192.168.33.11",
    "192.168.33.12",
    "192.168.33.13",
    "localhost",
    "127.0.0.1"
  ],
  "key": {
    "algo": "rsa",
    "size": 4096
  }
}
EOF

cat server-csr.json;

cfssl gencert -initca ca-csr.json | cfssljson -bare ca
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=config.json -profile=kubernetes server-csr.json | cfssljson -bare server
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=config.json -profile=client client-csr.json | cfssljson -bare client
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=config.json -profile=peer peer-csr.json | cfssljson -bare peer
```

```text
tree
.
|-- ca-csr.json
|-- ca-key.pem
|-- ca.csr
|-- ca.pem
|-- client-csr.json
|-- client-key.pem
|-- client.csr
|-- client.pem
|-- config.json
|-- peer-csr.json
|-- peer-key.pem
|-- peer.csr
|-- peer.pem
|-- server-csr.json
|-- server-key.pem
|-- server.csr
`-- server.pem

0 directories, 17 files
```

## Node installation

* Copy all cert to /etc/etcd/certs/
* Write a template config file, etc:

```yaml
name: 'default'
data-dir:
wal-dir:
snapshot-count: 10000
heartbeat-interval: 100
election-timeout: 1000
quota-backend-bytes: 0
listen-peer-urls: http://localhost:2380
listen-client-urls: http://localhost:2379
max-snapshots: 5
max-wals: 5
cors:
initial-advertise-peer-urls: http://localhost:2380
advertise-client-urls: http://localhost:2379
discovery:
discovery-fallback: 'proxy'
discovery-proxy:
discovery-srv:
initial-cluster:
initial-cluster-token: 'etcd-cluster'
initial-cluster-state: 'new'
strict-reconfig-check: false
enable-v2: true
enable-pprof: true
proxy: 'off'
proxy-failure-wait: 5000
proxy-refresh-interval: 30000
proxy-dial-timeout: 1000
proxy-write-timeout: 5000
proxy-read-timeout: 0
client-transport-security:
  cert-file: /etc/etcd/certs/server.pem
  key-file: /etc/etcd/certs/server-key.pem
  client-cert-auth: true
  trusted-ca-file: /etc/etcd/certs/ca.pem
  auto-tls: false
peer-transport-security:
  cert-file: /etc/etcd/certs/peer.pem
  key-file: /etc/etcd/certs/peer-key.pem
  client-cert-auth: true
  trusted-ca-file: /etc/etcd/certs/ca.pem
  auto-tls: false
debug: false
logger: zap
log-outputs: [stderr]
force-new-cluster: false
auto-compaction-mode: periodic
auto-compaction-retention: "1"
```

分别在192.168.33.11、192.168.33.12、192.168.33.13执行

```shell
mkdir -p /opt/etcd-data
chmod 700 /opt/etcd-data
mkdir -p /etc/etcd
mkdir -p /etc/etcd/certs
cat > /lib/systemd/system/etcd.service <<EOF
[Unit]
Description=etcd
Documentation=https://github.com/coreos/etcd
After=network.target
After=network-online.target
Wants=network-online.target
[Service]
Type=notify
ExecStart=/usr/local/bin/etcd --config-file /etc/etcd/etcd.conf.yml
Restart=on-failure
RestartSec=10s
LimitNOFILE=65536
[Install]
WantedBy=multi-user.target
EOF

ip4=$(/sbin/ip -o -4 addr list eth1 | awk '{print $4}' |cut -d/ -f1 | head -n1);
tmpn=$(echo -n ${ip4} | awk -F "." '{print $NF}');
cluster="etcd-11=https://192.168.33.11:2380,etcd-12=https://192.168.33.12:2380,etcd-13=https://192.168.33.13:2380"
#ip4="192.168.33.11"
#name="etcd-11"
sudo cp /tmp/etcd.conf.yml /etc/etcd && \
sudo sed -i "s@^name:.*@name: 'etcd-${tmpn}'@g" /etc/etcd/etcd.conf.yml && \
sudo sed -i "s@^data-dir:.*@data-dir: /opt/etcd-data@g" /etc/etcd/etcd.conf.yml && \
sudo sed -i "s@^listen-peer-urls:.*@listen-peer-urls: https://${ip4}:2380@g" /etc/etcd/etcd.conf.yml && \
sudo sed -i "s@^listen-client-urls:.*@listen-client-urls: https://${ip4}:2379@g" /etc/etcd/etcd.conf.yml && \
sudo sed -i "s@^initial-advertise-peer-urls:.*@initial-advertise-peer-urls: https://${ip4}:2380@g" /etc/etcd/etcd.conf.yml && \
sudo sed -i "s@^advertise-client-urls:.*@advertise-client-urls: https://${ip4}:2379@g" /etc/etcd/etcd.conf.yml && \
sudo sed -i "s@^initial-cluster:.*@initial-cluster: ${cluster}@g" /etc/etcd/etcd.conf.yml && \
cat /etc/etcd/etcd.conf.yml && \
cat /lib/systemd/system/etcd.service

sudo systemctl daemon-reload
sudo systemctl enable etcd.service
sudo systemctl restart etcd.service
sudo systemctl status etcd.service -l
```

```shell
etcdctl endpoint health \
--endpoints "https://192.168.33.11:2379,https://192.168.33.12:2379,https://192.168.33.13:2379" \
--cacert=/etc/etcd/certs/ca.pem \
--cert=/etc/etcd/certs/client.pem \
--key=/etc/etcd/certs/client-key.pem \
--cluster=true

https://192.168.33.11:2379 is healthy: successfully committed proposal: took = 48.748634ms
https://192.168.33.12:2379 is healthy: successfully committed proposal: took = 49.391402ms
https://192.168.33.13:2379 is healthy: successfully committed proposal: took = 54.411539ms
```
