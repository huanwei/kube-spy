#!/bin/bash

docker rm spy-influxdb -f
#rm -rf $PWD/spy-influxdb
docker run --name=spy-influxdb -p 8086:8086 -v $PWD/spy-influxdb:/var/lib/influxdb -e INFLUXDB_DB=spydb -e INFLUXDB_USER=spy -e INFLUXDB_USER_PASSWORD=123456 -d influxdb:alpine
