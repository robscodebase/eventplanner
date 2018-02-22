#!/bin/bash
#sudo docker stop $(docker ps -a -q) && sudo docker rm $(docker ps -a -q)
sudo docker stop mysql && sudo docker rm mysql
sudo docker build -t mysql .
sudo docker run --net eventNet --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD="insecure" -d  mysql
sudo docker exec -ti mysql /bin/bash
