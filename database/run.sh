#!/bin/bash
docker stop mysql-event-planner && docker rm mysql-event-planner
docker run --net eventNet --name mysql-event-planner -e MYSQL_ROOT_PASSWORD=insecure -d mysql
docker exec -ti mysql-event-planner /bin/bash
