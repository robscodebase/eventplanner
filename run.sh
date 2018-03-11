#!/bin/bash
#sudo docker stop $(docker ps -a -q) && sudo docker rm $(docker ps -a -q)
#docker run --net eventNet --name mysql-event-planner -e MYSQL_ROOT_PASSWORD="insecure" -d mysql
docker rmi $(docker images -a | grep "<none>")
docker stop event-planner && sudo docker rm event-planner
docker build -t event-planner .
docker run -dp 8081:8081 --name event-planner --net eventNet event-planner
docker exec -ti event-planner /bin/bash
