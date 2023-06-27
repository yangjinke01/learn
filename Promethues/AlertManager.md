```yaml
global:
  scrape_interval: 15s
  external_labels:
    monitor: 'monitor'
rule_files:
  - rules.yml
alerting:
  alertmanagers:
    - static_configs:
        - targets: [ '10.0.111.146:9093' ]
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: [ 'artefact:9100','nexus:9100','confluence:9100' ]
```

rules.yml
```yaml
groups:
  - name: example
    rules:
      - alert: InstanceDown
        expr: up == 0
        for: 1m
```