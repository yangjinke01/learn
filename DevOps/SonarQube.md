# SonarQube

## Docker Install

没有使用外部数据库，页面会有提示，测试用

```shell
docker volume create --name sonarqube_data
docker volume create --name sonarqube_logs
docker volume create --name sonarqube_conf
docker volume create --name sonarqube_extensions

docker run -d --name sonarqube \
    -p 9000:9000 \
    -e SONAR_ES_BOOTSTRAP_CHECKS_DISABLE=true \
    -v sonarqube_data:/opt/sonarqube/data \
    -v sonarqube_conf:/opt/sonarqube/conf \
    -v sonarqube_extensions:/opt/sonarqube/extensions \
    -v sonarqube_logs:/opt/sonarqube/logs \
    sonarqube:9.6-enterprise
    
# admin/admin
```

## Docker Compose (推荐)

```yaml
# vim docker-compose.yml

version: "3.9"

services:
  sonarqube:
    image: sonarqube:lts-community
    restart: always
    container_name: sonarqube
    depends_on:
      - postgres
    ports:
      - "9000:9000"
    environment:
      SONAR_JDBC_URL: jdbc:postgresql://postgres:5432/sonar
      SONAR_JDBC_USERNAME: sonar
      SONAR_JDBC_PASSWORD: sonar
    volumes:
      - SonarQube_data:/opt/sonarqube/data
      - SonarQube_extensions:/opt/sonarqube/extensions
      - SonarQube_logs:/opt/sonarqube/logs
      - SonarQube_conf:/opt/sonarqube/conf
    networks:
      - sonar
  postgres:
    image: postgres:12
    restart: always
    container_name: postgres
    ports:
      - "5432:5432"
    environment:
      TZ: Asia/Shanghai
      POSTGRES_USER: sonar
      POSTGRES_PASSWORD: sonar
      POSTGRES_DB: sonar
    volumes:
      - postgresql:/var/lib/postgresql
      - postgresql_data:/var/lib/postgresql/data
    networks:
      - sonar

volumes:
  SonarQube_data:
  SonarQube_extensions:
  SonarQube_logs:
  SonarQube_conf:
  postgresql:
  postgresql_data:

networks:
  sonar:
    driver: bridge

# docker-compose up -d
```

## LDAP 配置

```text
# /opt/sonarqube/conf/sonar.properties
sonar.security.realm=LDAP
ldap.url=ldap://10.0.31.210:10389
ldap.bindDn=uid=admin,ou=system
ldap.bindPassword=secret
ldap.user.baseDn=dc=ecloud,dc=com
```

## 集成GitLab

```yaml
stages:
  - SCAN

sonarqube-check:
  stage: SCAN
  image:
    name: sonarsource/sonar-scanner-cli:4
  variables:
    SONAR_HOST_URL: 'http://10.0.31.210:9000'
    SONAR_TOKEN: "sqa_16e48733d3e0c95eb09ffdd0f6ab069b481b3328"
  script:
    - echo
      SONAR_HOST_URL=${SONAR_HOST_URL} SONAR_LOGIN=${SONAR_TOKEN} sonar-scanner
      -Dsonar.qualitygate.wait=true
      -Dsonar.projectKey=${CI_PROJECT_NAME}_${CI_COMMIT_REF_NAME}
      -Dsonar.projectName=${CI_PROJECT_NAME}_${CI_COMMIT_REF_NAME}
      -Dsonar.projectVersion=${CI_COMMIT_SHORT_SHA}
      -Dsonar.sourceEncoding=UTF-8
    - SONAR_HOST_URL=${SONAR_HOST_URL} SONAR_LOGIN=${SONAR_TOKEN} sonar-scanner
      -Dsonar.qualitygate.wait=true
      -Dsonar.projectKey=${CI_PROJECT_NAME}_${CI_COMMIT_REF_NAME}
      -Dsonar.projectName=${CI_PROJECT_NAME}_${CI_COMMIT_REF_NAME}
      -Dsonar.projectVersion=${CI_COMMIT_SHORT_SHA}
      -Dsonar.sourceEncoding=UTF-8
  allow_failure: true
  only:
    - main
  tags:
    - docker
```