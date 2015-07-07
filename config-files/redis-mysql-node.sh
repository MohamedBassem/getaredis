#!/bin/bash

MYSQL_PASSOWRD=
REDIS_PASSWORD=

#Installing MySQL
apt-get update
echo mysql-server mysql-server/root_password password $MYSQL_PASSOWRD | sudo debconf-set-selections
echo mysql-server mysql-server/root_password_again password $MYSQL_PASSOWRD | sudo debconf-set-selections
apt-get install -y mysql-server mysql-client

#Installing Redis
add-apt-repository -y ppa:rwky/redis
apt-get update
apt-get install -y redis-server
echo "requirepass $REDIS_PASSWORD" >> /etc/redis/redis.conf
service redis-server restart

#Installing htop
apt-get install -y htop
