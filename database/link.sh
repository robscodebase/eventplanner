#!/bin/bash
docker run --net eventNet -it --link mysql-event-planner:mysql --rm mysql sh -c 'exec mysql -h"mysql-event-planner" -P"3306" -u root -p"$MYSQL_ENV_MYSQL_ROOT_PASSWORD"'
