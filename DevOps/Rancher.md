# Rancher

## installation

```shell
docker run -d --restart=unless-stopped \
  -p 80:80 -p 443:443 \
  -v /opt/rancher:/var/lib/rancher \
  --privileged \
  rancher/rancher:latest
```

kubectl delete crd backingimagedatasources.longhorn.io
kubectl delete crd backingimagemanagers.longhorn.io
kubectl delete crd backingimages.longhorn.io
kubectl delete crd backups.longhorn.io
kubectl delete crd backuptargets.longhorn.io
kubectl delete crd backupvolumes.longhorn.io
kubectl delete crd engineimages.longhorn.io
kubectl delete crd engines.longhorn.io
kubectl delete crd instancemanagers.longhorn.io
kubectl delete crd nodes.longhorn.io
kubectl delete crd orphans.longhorn.io
kubectl delete crd recurringjobs.longhorn.io
kubectl delete crd replicas.longhorn.io
kubectl delete crd settings.longhorn.io
kubectl delete crd sharemanagers.longhorn.io
kubectl delete crd snapshots.longhorn.io
kubectl delete crd supportbundles.longhorn.io
kubectl delete crd systembackups.longhorn.io
kubectl delete crd systemrestores.longhorn.io
kubectl delete crd volumes.longhorn.io
