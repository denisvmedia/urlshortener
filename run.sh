#!/bin/sh

docker-compose up -d mariadb
sleep 3
docker-compose up -d shortener
