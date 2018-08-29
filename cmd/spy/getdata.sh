#!/usr/bin/env bash
# arg 1: search time, example: 5m(5 minutes)
# Get db address
list=$(kubectl get pod -o wide|grep influxdb-spy)
# Split string
array=(${list// / })
# Query
curl -G "http://${array[5]}:8086/query?pretty=true" --data-urlencode "db=spy" --data-urlencode "q=SELECT * FROM response WHERE time > now() - $1 LIMIT 4"
curl -G "http://${array[5]}:8086/query?pretty=true" --data-urlencode "db=spy" --data-urlencode "q=SELECT * FROM ping WHERE time > now() - $1 "