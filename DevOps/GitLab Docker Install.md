# GitLab Docker Install

## Install GitLab

```shell
export GITLAB_HOME=/etc/gitlab

sudo docker run --detach \
  --hostname gitlab.jack.com \
  --publish 443:443 --publish 80:80 --publish 222:22 \
  --name gitlab \
  --restart always \
  --volume $GITLAB_HOME/config:/etc/gitlab:Z \
  --volume $GITLAB_HOME/logs:/var/log/gitlab:Z \
  --volume $GITLAB_HOME/data:/var/opt/gitlab:Z \
  --shm-size 256m \
  gitlab/gitlab-ce:15.2.2-ce.0

# get root password
sudo docker exec -it gitlab grep 'Password:' /etc/gitlab/initial_root_password
```

## Install GitLab runner

```shell
# Download the binary for your system
sudo curl -L --output /usr/local/bin/gitlab-runner https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64

# Give it permission to execute
sudo chmod +x /usr/local/bin/gitlab-runner

# Create a GitLab Runner user
sudo useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash

# Install and run as a service
sudo gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner
sudo gitlab-runner start

```

## Register Runner

```shell
# runner 得能解析 gitlab.ecloud.com
gitlab-runner register --non-interactive \
    --url http://gitlab.ecloud.com/ \
    --registration-token p2P4_-ZyLpJnzAF7AxSA \
    --executor "docker" \
    --description "docker" \
    --tag-list "docker" \
    --docker-extra-hosts "gitlab.ecloud.com:10.0.31.210" \
    --docker-image "centos:7"
```
