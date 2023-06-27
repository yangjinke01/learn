# Rook Deployment

## download image from k8s.gcr.io script

```shell
cat << EOF > pullK8sImage.sh
#!/bin/sh

k8s_img=\$1
mirror_img=\$(echo \${k8s_img}|
        sed 's/k8s\.gcr\.io/anjia0532\/google-containers/g;s/gcr\.io/anjia0532/g;s/\//\./g;s/ /\n/g;s/anjia0532\./anjia0532\//g' |
        uniq)
echo \$mirror_img
sudo docker pull \${mirror_img}
sudo docker tag \${mirror_img} \${k8s_img}
sudo docker rmi \${mirror_img}
EOF

chmod +x pullK8sImage.sh

./pullK8sImage.sh k8s.gcr.io/sig-storage/csi-attacher:v3.4.0
./pullK8sImage.sh k8s.gcr.io/sig-storage/csi-node-driver-registrar:v2.5.0
./pullK8sImage.sh k8s.gcr.io/sig-storage/csi-provisioner:v3.1.0
./pullK8sImage.sh k8s.gcr.io/sig-storage/csi-resizer:v1.4.0
./pullK8sImage.sh k8s.gcr.io/sig-storage/csi-snapshotter:v5.0.1
./pullK8sImage.sh k8s.gcr.io/sig-storage/nfsplugin:v3.1.0

docker save -o k8s.tgz k8s.gcr.io/sig-storage/csi-attacher:v3.4.0 k8s.gcr.io/sig-storage/csi-node-driver-registrar:v2.5.0 k8s.gcr.io/sig-storage/csi-provisioner:v3.1.0 k8s.gcr.io/sig-storage/csi-resizer:v1.4.0 k8s.gcr.io/sig-storage/csi-snapshotter:v5.0.1 k8s.gcr.io/sig-storage/nfsplugin:v3.1.0
```

## load images in every k8s node

```shell
ctr -n=k8s.io image import k8s.tgz
```

## deploy Rook on k8s

官方文档 <https://rook.io/docs/rook/v1.9/Getting-Started/quickstart/>
每个节点的内核需要在4.17以上

```shell
export HTTPS_PROXY=socks5://10.0.10.113:7070

git clone --single-branch --branch v1.9.4 https://github.com/rook/rook.git
cd rook/deploy/examples
kubectl create -f crds.yaml -f common.yaml -f operator.yaml
kubectl create -f cluster.yaml

# all pods status is running
kubectl -n rook-ceph get pod
```

## Toolbox

<https://rook.io/docs/rook/v1.9/Troubleshooting/ceph-toolbox/>

```shell
kubectl create -f deploy/examples/toolbox.yaml
kubectl -n rook-ceph rollout status deploy/rook-ceph-tools
kubectl -n rook-ceph exec -it deploy/rook-ceph-tools -- bash

ceph status
ceph osd status
ceph df
rados df
```

## Dashboard

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: rook-ceph-mgr-dashboard
  namespace: rook-ceph
  annotations:
    kubernetes.io/ingress.class: "nginx"
    kubernetes.io/tls-acme: "true"
    nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
    nginx.ingress.kubernetes.io/server-snippet: |
      proxy_ssl_verify off;
spec:
  tls:
   - hosts:
     - rook-ceph.example.com
     secretName: rook-ceph.example.com
  rules:
  - host: rook-ceph.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: rook-ceph-mgr-dashboard
            port:
              name: https-dashboard
```
