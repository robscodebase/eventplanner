#!/bin/bash
#sudo docker stop $(docker ps -a -q) && sudo docker rm $(docker ps -a -q)
sudo docker stop mysqldb && sudo docker rm mysqldb
sudo docker build -t mysqldb .
sudo docker run --name mysqldb -e MYSQL_ROOT_PASSWORD="insecure" -d  mysqldb
sudo docker exec -ti mysqldb /bin/bash
