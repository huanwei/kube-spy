#!/usr/bin/env bash

for((i=${1};i>=1;i--));
do
appName="http-test"${i}".yaml"
currentServiceName="http-test-service"${i}".yaml"
kubectl apply -f ${appName}
kubectl apply -f ${currentServiceName}
done
