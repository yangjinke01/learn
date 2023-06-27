# Ubuntu install Mysql

## setting root password

```shell
apt update
apt install mysql-server
systemctl start mysql.service

mysql
mysql> ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '******';
mysql> exit
```

##  Creating a Dedicated MySQL User and Granting Privileges

```shell
mysql -u root -p
mysql> CREATE USER 'jtck'@'%' IDENTIFIED BY 'Yjk_13525748624';
```

```shell
CREATE DATABASE jtck CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL PRIVILEGES ON jtck.* TO 'jtck'@'%';
FLUSH PRIVILEGES;
```

```shell
vim /etc/mysql/mysql.conf.d/mysqld.cnf
```

