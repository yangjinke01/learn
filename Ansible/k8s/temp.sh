#!/bin/bash
sed -i s@'sandbox_image = "k8s.gcr.io/pause:3.6"'@'sandbox_image = "registry.aliyuncs.com/google_containers/pause:3.7"'@ /etc/containerd/config.toml
