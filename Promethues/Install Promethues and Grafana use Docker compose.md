# Install Promethues and Grafana use Docker compose

## Install Prometheus

### Config of Prometheus

```shell
mkdir /etc/prometheus
vim /etc/prometheus/prometheus.yml
```

```yaml
global:
  scrape_interval: 15s
  external_labels:
    monitor: 'codelab-monitor'
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```

### Docker Compose up file docker-compose.yml

```yaml
version: '3'

services:
  prometheus:
    image: prom/prometheus
    ports:
      - 9090:9090
    volumes:
      - /etc/prometheus:/etc/prometheus
      - prometheus-data:/prometheus
    command: --web.enable-lifecycle  --config.file=/etc/prometheus/prometheus.yml

volumes:
  prometheus-data:
```

```shell
docker compose up -d
```

command: --web.enable-lifecycle --config.file=/etc/prometheus/prometheus.yml is optional. If you use --web.enable-lifecycle you can reload configuration files (e.g. rules) without restarting Prometheus:

```shell
curl -X POST http://localhost:9090/-/reload
```

## Docker compose install grafana

docker-compose.yml

```yaml
version: '3'

services:
  grafana:
      image: grafana/grafana
      ports:
        - 3000:3000
      restart: unless-stopped
      volumes:
        - grafana-config:/etc/grafana
        - grafana-data:/var/lib/grafana

volumes:
  grafana-config:
  grafana-data:
```

```shell
docker-compose up -d
```
