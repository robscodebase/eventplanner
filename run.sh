#!/bin/bash
#docker stop $(docker ps -a -q) && docker rm $(docker ps -a -q)
#docker rmi $(docker images -a | grep "event")
#docker rmi $(docker images -a | grep "<none>")
#docker stop event-planner && sudo docker rm event-planner
#docker run --net eventNet --name mysql-event-planner -e MYSQL_ROOT_PASSWORD="insecure" -d mysql
#docker build -t event-planner .
#docker run -dp 8081:8081 --name event-planner --net eventNet event-planner
#docker exec -ti event-planner /bin/bash

#
#
#Test
docker stop $(docker ps -a -q) && docker rm $(docker ps -a -q)
docker rmi $(docker images -a | grep "event")
docker rmi $(docker images -a | grep "<none>")
docker run --name mysql-event-planner-test -e MYSQL_ROOT_PASSWORD="insecure" -d mysql
docker build -t event-planner-test .
docker run -dp 8081:8081 --name event-planner-test event-planner-test
#docker exec -ti event-planner-test /bin/bash
