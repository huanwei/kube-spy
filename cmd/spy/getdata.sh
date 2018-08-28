#!/usr/bin/env bash
# Get db address
list=$(kubectl get pod -o wide|grep influxdb-spy)
# Split string
array=(${list// / })
# Query
curl -G "http://${array[5]}:8086/query?pretty=true" --data-urlencode "db=spy" --data-urlencode "q=SELECT * FROM response LIMIT 4"
curl -G "http://${array[5]}:8086/query?pretty=true" --data-urlencode "db=spy" --data-urlencode "q=SELECT * FROM ping 9"