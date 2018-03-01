#!/bin/bash
#sudo docker stop $(docker ps -a -q) && sudo docker rm $(docker ps -a -q)
sudo docker stop event-planner && sudo docker rm event-planner
sudo docker build -t event-planner .
sudo docker run -dp 8081:8081 --name event-planner --net eventNet event-planner
sudo docker exec -ti event-planner /bin/bash
