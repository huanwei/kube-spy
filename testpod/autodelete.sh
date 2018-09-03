#!/usr/bin/env bash

for((i=1;i<=${1};i++));
do
appName="http-test"${i}
currentServiceName="http-test-service"${i}

kubectl delete service ${currentServiceName}
kubectl delete deployment ${appName}
done