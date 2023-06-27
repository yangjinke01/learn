# Install Node.js via binary archive on Linux

```shell
mkdir -p /usr/local/lib/nodejs
tar -xf node-v18.12.1-linux-x64.tar.xz -C /usr/local/lib/nodejs/

vim /etc/profile
# export PATH=/usr/local/lib/nodejs/node-v18.12.1-linux-x64/bin:$PATH

. /etc/profile
npm version
```

