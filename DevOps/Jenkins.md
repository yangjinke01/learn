# Jenkins

## 安装

```shell
cat > /etc/yum.repos.d/jenkins.repo << EOF
[jenkins]
name=Jenkins-stable
baseurl=http://pkg.jenkins.io/redhat-stable
gpgcheck=0
EOF

yum install -y jenkins java-11-openjdk
systemctl start jenkins
```



## 修改为国内源

```shell
cd  /var/lib/jenkins/
find  -name "default.json" 
#/var/lib/jenkins/updates/default.json

sed -i 's/http:\/\/updates.jenkins-ci.org\/download/https:\/\/mirrors.tuna.tsinghua.edu.cn\/jenkins/g' /var/lib/jenkins/updates/default.json && sed -i 's/http:\/\/www.google.com/https:\/\/www.baidu.com/g' /var/lib/jenkins/updates/default.json
systemctl restart jenkins
```

