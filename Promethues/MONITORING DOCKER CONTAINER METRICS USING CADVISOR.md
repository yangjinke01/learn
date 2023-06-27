# MONITORING DOCKER CONTAINER METRICS USING CADVISOR

## Prometheus configuration

```yaml
scrape_configs:
- job_name: cadvisor
  scrape_interval: 5s
  static_configs:
  - targets:
    - cadvisor:8080
```

## Pull cadvisor image

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

./pullK8sImage.sh gcr.io/cadvisor/cadvisor:latest
```

## Docker Compose configuration file docker-compose.yml

```yaml
version: '3'
services:
  cadvisor:
    image: gcr.io/cadvisor/cadvisor:latest
    container_name: cadvisor
    ports:
    - 8080:8080
    volumes:
    - /:/rootfs:ro
    - /var/run:/var/run:rw
    - /sys:/sys:ro
    - /var/lib/docker/:/var/lib/docker:ro
```

```shell
docker-compose up -d
```
