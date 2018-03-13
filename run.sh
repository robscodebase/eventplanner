#!/bin/bash
docker build -t event-planner .
docker run -dp 8081:8081 --name event-planner --net eventNet event-planner
