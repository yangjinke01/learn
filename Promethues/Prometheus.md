# Prometheus

## Introduction

Prometheus is an open-source monitoring and alerting toolkit originally developed at SoundCloud. It is designed to help you collect, store, query, and visualize time-series data, and is particularly well-suited for monitoring and alerting in dynamic, cloud-native environments. Prometheus is part of the Cloud Native Computing Foundation (CNCF) and is widely used in the DevOps and cloud-native communities.

Here are some key features and components of Prometheus:

1. **Data Collection**: Prometheus collects time-series data through a pull model. It scrapes metrics from instrumented applications and services at regular intervals. These metrics can be in the form of numerical data, text, or custom data types.

2. **Data Storage**: Collected data is stored locally in a time-series database. Prometheus uses a specialized storage format to efficiently store and query time-series data.

3. **Query Language**: Prometheus provides a powerful query language called PromQL (Prometheus Query Language) that allows users to perform complex queries and create custom alerts based on the collected data.

4. **Alerting**: Prometheus can trigger alerts based on user-defined alerting rules. When a metric exceeds a specified threshold, Prometheus can send notifications to various alerting channels, such as email, Slack, or other communication platforms.

5. **Visualization**: While Prometheus itself is not a visualization tool, it can be integrated with popular visualization tools like Grafana. Grafana is often used alongside Prometheus to create interactive and customizable dashboards for monitoring data.

6. **Service Discovery**: Prometheus supports service discovery mechanisms, making it suitable for dynamic environments like Kubernetes. It can automatically discover and scrape metrics from newly deployed services.

7. **Exporters**: Prometheus exporters are small software components that help bridge the gap between Prometheus and various third-party systems. Exporters expose metrics in a format Prometheus can understand.

8. **Reliability**: Prometheus is designed with reliability in mind. It offers local storage, meaning that if the network or remote systems fail, Prometheus can continue to scrape and store metrics.

Prometheus is widely used in container orchestration platforms like Kubernetes and is well-suited for monitoring microservices and other cloud-native applications. It provides valuable insights into the performance and health of your systems, making it easier to detect and respond to issues promptly.

To use Prometheus effectively, you'll typically set up Prometheus servers, configure them to scrape data from your applications and services, define alerting rules, and visualize the data using dashboards. Integrating it with Grafana is a common choice for creating rich and interactive monitoring dashboards.

## types of metrics

Prometheus supports several types of metrics that can be collected and used for monitoring and alerting. These metric types are important for understanding and querying the data effectively in Prometheus. The four primary metric types in Prometheus are:

1. **Counter**: Counters are used to represent values that can only increase over time, such as the number of requests processed or the total number of events. Counters are always non-negative and are typically used to measure the rate of events over time. When querying a counter, Prometheus can calculate the rate of change over a specified time period, which is useful for understanding trends and performance.

   Example: `http_requests_total` - the total number of HTTP requests received by a web server.

2. **Gauge**: Gauges are used for representing values that can go up and down, such as CPU usage, memory utilization, or temperature. Gauges are often used for instantaneous or point-in-time measurements. When you query a gauge, you get the current value.

   Example: `cpu_usage` - the current CPU usage percentage.

3. **Histogram**: Histograms are used to measure the distribution of values over time. They divide the values into buckets or bins, allowing you to observe the distribution, percentiles, and other statistical measures of a metric. Histograms are valuable for understanding the spread of values, especially in scenarios where you want to assess latency or response times.

   Example: `http_request_duration_seconds` - a histogram of HTTP request durations.

4. **Summary**: Summaries are similar to histograms but are more focused on calculating percentiles and quantiles of observed values. They provide a way to compute median, 90th percentile, 99th percentile, and other quantiles of a metric.

   Example: `api_response_time_seconds` - a summary of API response times.

These metric types are crucial for effective monitoring and alerting because they allow you to capture a wide range of information about your systems, from counts and instantaneous values to distributions and percentiles. By using these metric types, you can gain insights into the behavior and performance of your applications and infrastructure.

Prometheus's query language, PromQL, is specifically designed to work with these metric types, enabling you to perform various calculations and transformations on the collected data. You can use PromQL to create custom alerts, build dashboards, and gain a deeper understanding of your systems' health and performance.
