#!/bin/sh

docker-compose up -d mariadb
sleep 5
docker-compose run shortener_init
