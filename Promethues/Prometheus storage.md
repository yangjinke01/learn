# Prometheus storage: technical terms for humans

Many technical terms could be used when referring to Prometheus storage — either local storage or remote storage. New users could be unfamiliar with these terms, which could result in misunderstandings. Let’s explain the most commonly used technical terms in this article.

## Time Series

A time series is a series of (timestamp, value) pairs sorted by timestamp. The number of pairs per each time series can be arbitrary — from one to hundreds of billions. Timestamps have millisecond precision, while values are 64-bit floating point numbers. Each time series has a name. For instance:

* node_cpu_seconds_total — the total number of CPU seconds used
* node_filesystem_free_bytes — free space on filesystem mount point
* go_memstats_sys_bytes — the amount of memory used by Go app

Additionally to name each time series can have arbitrary number of label=”value” labels. For instance:

* go_memstats_sys_bytes{instance=”foobar:1234",job=”node_exporter”}
* prometheus_http_requests_total{handler=”/api/v1/query_range”}

Each time series is uniquely identified by its name plus a set of labels. For example, these all are distinct time series:

* temperature{city=”SF”}
* temperature{city=”SF”, unit=”Celsius”}

## Data Point or Sample

Each (timestamp, value) pair from any time series is called a data point or sample.

## Metric Types

The following metric types are supported by Prometheus:

* Gauge — a value, which can go up and down at any given time. For example, temperature or memory usage.
* Counter — a value, which starts from 0 and never goes down. For example, requests count or distance traveled. There is one exception, which is called counter reset, when counter resets to 0 if the service exposing the metric is restarted.
* Summary — maintains a set of pre-configured percentiles for the value.
* Histogram — maintains a set of counters (aka buckets) for different value ranges. The set of buckets can be pre-configured or can be static. See this article for details.

## High Cardinality

The number of unique time series stored in the TSDB is called cardinality. “High cardinality” means high number of series. Prometheus tsdb storage is optimized for working with relatively low number of time series — see [these slides from PromCon 2019](https://promcon.io/2019-munich/slides/containing-your-cardinality.pdf). It could start working slowly or could use high amounts of RAM, CPU or disk IO during ingestion, querying or [compaction](https://prometheus.io/docs/prometheus/latest/storage/#compaction) when high number of time series is stored in it. High cardinality for Prometheus starts from a few millions of time series. Read [this article about how to avoid high cardinality in Prometheus](https://www.robustperception.io/cardinality-is-key).

Prometheus exposes the information about high cardinality time series at `/status` page starting from [v2.14.0](https://github.com/prometheus/prometheus/releases) — see [this PR](https://github.com/prometheus/prometheus/pull/6125) for details.

## Active time series

A time series is considered active if Prometheus scraped new data for it recently. Prometheus provides `prometheus_tsdb_head_series` metric, which shows the number of active time series. Prometheus holds recently added samples for active time series in RAM, so its RAM usage highly depends on the number of active time series. The number of active time series is connected to cardinality via churn rate. See [this article on how to estimate RAM usage for Prometheus from the number of active time series and ingestion rate](https://www.robustperception.io/how-much-ram-does-prometheus-2-x-need-for-cardinality-and-ingestion).

## Churn rate

Active time series becomes inactive if it stops receiving new samples. New time series can substitute old time series on label value change. For instance, `pod_name` label value may change for big number of time series after each deployment in Kubernetes. The rate at which old time series are substituted by new time series is called churn rate. Obviously, high churn rate increases the cardinality. This may lead to performance degradation, high RAM usage and out of memory errors aka OOM.

Prometheus exposes `prometheus_tsdb_head_series_created_total` metric, which could be used for estimating the churn rate using the following PromQL query:

```text
rate(prometheus_tsdb_head_series_created_total[5m])
```

Starting from v2.10 Prometheus [exposes](https://www.robustperception.io/finding-churning-targets-in-prometheus-with-scrape_series_added) per-target `scrape_series_added` metric, which can be used for determining the source of high series churn rate:

```text
sum(sum_over_time(scrape_series_added[5m])) by (job)
```

See also [this article](https://fabxc.org/tsdb/), which explains churn rate in Prometheus with more details.

## Scrape interval

Prometheus scrapes targets with the configured interval, which is named scrape interval. By default it equals to 1 minute for all the scrape targets. It can be overridden via `global->scrape_interval` option in [Prometheus config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#configuration-file). While it is possible to set distinct scrape interval per each target, this isn’t recommended. See [this article](https://www.robustperception.io/keep-it-simple-scrape_interval-id) for details.

Lower scrape interval results in higher ingestion rate and in higher RAM usage for Prometheus, since more data points must be kept in RAM before they are flushed to disk.

## Retention

Prometheus provides `--storage.tsdb.retention.time` command-line flag for configuring the lifetime for the stored data — see [these docs](https://prometheus.io/docs/prometheus/latest/storage/#operational-aspects) for more info. The data outside the retention is automatically deleted. By default the retention is configured to [15 days](https://github.com/prometheus/prometheus/blob/0ea3a2218d3a71d7a721c078efa2919175beb7a4/cmd/prometheus/main.go#L75). The amounts of data stored on disk depends on retention — higher retention means more data on disk.

The lowest supported retention in Prometheus is 2 hours (2h). Such a retention could be useful when configuring [remote storage for Prometheus](https://github.com/VictoriaMetrics/VictoriaMetrics/blob/master/README.md#prometheus-setup). In this case Prometheus simultaneously replicates all the scraped data into local storage and all the configured remote storage backends. This means that the retention for local storage can be minimal, since all the data is already replicated to remote storage and the [remote storage can be used for querying from Grafana](https://github.com/VictoriaMetrics/VictoriaMetrics/blob/master/README.md#grafana-setup) and any other clients with [Prometheus querying API](https://prometheus.io/docs/prometheus/latest/querying/api/) support.

Note that the configured retention must cover time ranges for [alerting](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/) and [recording rules](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/).

## Relabeling

Prometheus supports relabeling, which can be used for manipulating per-metric labels, filtering scrape targets and samples. This is quite complex but very powerful mechanism, which is frequently used with service discovery for Kubernetes, Amazon EC2, Google Compute Engine, etc. See [this article](https://valyala.medium.com/how-to-use-relabeling-in-prometheus-and-victoriametrics-8b90fc22c4b2) for details.

## Remote Storage

By default Prometheus stores data to local tsdb, but it can be [configured to replicate the scraped data to remote storage backends](https://github.com/VictoriaMetrics/VictoriaMetrics/blob/master/README.md#prometheus-setup). This can be useful in the following cases:

* Collecting data from many Prometheus instances to a single remote storage, so all the data could be queried and analyzed. This is sometimes called `global query view`.
* Storing long-term data into remote storage, so local tsdb in Prometheus could be configured with low retention in order to occupy low amounts of disk space.
* Overcoming scalability issues for Prometheus, which cannot automatically scale to multiple nodes. Certain remote storage solutions such as VictoriaMetrics can scale both [vertically](https://medium.com/@valyala/measuring-vertical-scalability-for-time-series-databases-in-google-cloud-92550d78d8ae) (i.e. on a single computer) and [horizontally](https://github.com/VictoriaMetrics/VictoriaMetrics/blob/cluster/README.md) (i.e. clustering over multiple computers).
* Running Prometheus in K8S cluster with ephemeral storage volumes, which can disappear after pod restart (aka stateless mode). This is safe, since all the scraped data is immediately replicated to the configured remote storage backends.

## Conclusion

I hope this article helps clearing the meaning for popular tech terms used in the context of Prometheus storage.

P.S. I’m the author of [VictoriaMetrics](https://github.com/VictoriaMetrics/VictoriaMetrics/) — open source cost effective remote storage for Prometheus with easy setup and operation. I’d recommend taking a look at it if you use Prometheus at work. And join our [Slack chat](http://slack.victoriametrics.com/).
