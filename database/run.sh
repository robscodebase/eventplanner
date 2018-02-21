#!/bin/bash
#sudo docker stop $(docker ps -a -q) && sudo docker rm $(docker ps -a -q)
sudo docker stop mysqldb && sudo docker rm mysqldb
sudo docker build -t mysqldb .
sudo docker run --name mysqldb -p 3306:3306 -e MYSQL_ROOT_PASSWORD="insecure" -e MYSQL_ROOT_USERNAME="mysqldb" -d  mysqldb
sudo docker exec -ti mysqldb /bin/bash
