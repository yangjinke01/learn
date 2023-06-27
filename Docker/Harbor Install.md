# Harbor Install

## Prerequisites

### Software

The following table lists the software versions that must be installed on the target host.

| Software       | Version                       | Description                                                  |
| :------------- | :---------------------------- | :----------------------------------------------------------- |
| Docker engine  | Version 17.06.0-ce+ or higher | For installation instructions, see [Docker Engine documentation](https://docs.docker.com/engine/installation/) |
| Docker Compose | Version 1.18.0 or higher      | For installation instructions, see [Docker Compose documentation](https://docs.docker.com/compose/install/) |
| Openssl        | Latest is preferred           | Used to generate certificate and keys for Harbor             |

### Network ports

Harbor requires that the following ports be open on the target host.

| Port | Protocol | Description                                                  |
| :--- | :------- | :----------------------------------------------------------- |
| 443  | HTTPS    | Harbor portal and core API accept HTTPS requests on this port. You can change this port in the configuration file. |
| 4443 | HTTPS    | Connections to the Docker Content Trust service for Harbor. Only required if Notary is enabled. You can change this port in the configuration file. |
| 80   | HTTP     | Harbor portal and core API accept HTTP requests on this port. You can change this port in the configuration file. |

## Download And Install

```shell
# https://github.com/goharbor/harbor/releases
tar -xf harbor-offline-installer-v2.5.2.tgz -C /opt/
cd /opt/

# https
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -sha512 -days 3650 \
 -subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=yourdomain.com" \
 -key ca.key \
 -out ca.crt
mkdir cert
mv ca.key ca.crt cert

cp harbor.yml.tmpl harbor.yml
vim harbor.yml
# hostname: 10.0.31.220
#   certificate: /opt/harbor/cert/ca.crt
#   private_key: /opt/harbor/cert/ca.key

./prepare
./install.sh
```

docker pull 10.0.31.220:8443/docker-official/library/nginx@sha256:3536d368b898eef291fb1f6d184a95f8bc1a6f863c48457395aab859fda354d1