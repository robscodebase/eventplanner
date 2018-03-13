#!/bin/bash
docker stop $(docker ps -a -q) && docker rm $(docker ps -a -q)
docker rmi $(docker images -a | grep "event")
docker rmi $(docker images -a | grep "<none>")
docker run --name mysql-event-planner-test -e MYSQL_ROOT_PASSWORD="insecure" -d mysql
docker build -t event-planner-test .
docker run -dp 8081:8081 --name event-planner-test event-planner-test
