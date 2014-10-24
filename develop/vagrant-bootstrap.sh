#!/usr/bin/env bash

MYSQL_PASSWORD="pwd"
MYSQL_DBNAME="tima"

apt-get update

echo "mysql-server-5.5 mysql-server/root_password password $MYSQL_PASSWORD" | debconf-set-selections
echo "mysql-server-5.5 mysql-server/root_password_again password $MYSQL_PASSWORD" | debconf-set-selections
apt-get -y install mysql-server > /dev/null 2>&1

sed -i 's/^bind-address.*/bind-address = 0.0.0.0/' /etc/mysql/my.cnf

echo "create database $MYSQL_DBNAME;" | mysql -uroot -p$MYSQL_PASSWORD
echo "grant all privileges on *.* to 'root'@'%' identified by 'pwd' with grant option;" | mysql -uroot -p$MYSQL_PASSWORD

service mysql restart
