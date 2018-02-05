#!/bin/bash
sudo docker stop $(docker ps -a -q) && sudo docker rm $(docker ps -a -q)
sudo docker build -t event-planner .
sudo docker run -dp 8080:8080 --name event-planner event-planner
sudo docker exec -ti event-planner /bin/bash
