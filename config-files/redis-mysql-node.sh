#!/bin/bash

MYSQL_PASSWORD=
REDIS_PASSWORD=
DATABASE_NAME=getaredis

#Installing MySQL
apt-get update
echo mysql-server mysql-server/root_password password $MYSQL_PASSWORD | sudo debconf-set-selections
echo mysql-server mysql-server/root_password_again password $MYSQL_PASSWORD | sudo debconf-set-selections
apt-get install -y mysql-server mysql-client

#Allowing MySQL remote access
sed -i 's/\(bind-address.*=.*\)127.0.0.1/\1 0.0.0.0/g' /etc/mysql/my.cnf
echo "CREATE DATABASE $DATABASE_NAME; GRANT ALL ON $DATABASE_NAME.* TO root@'%' IDENTIFIED BY '$MYSQL_PASSWORD'" > /tmp/mysqltmp
mysql -p$MYSQL_PASSWORD < /tmp/mysqltmp
rm /tmp/mysqltmp
service mysql restart

#Installing Redis
add-apt-repository -y ppa:rwky/redis
apt-get update
apt-get install -y redis-server
echo "requirepass $REDIS_PASSWORD" >> /etc/redis/redis.conf
service redis-server restart

#Installing htop
apt-get install -y htop
