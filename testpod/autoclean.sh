#!/usr/bin/env bash

for((i=1;i<=${1};i++));
do
appName="http-test"${i}".yaml"
currentServiceName="http-test-service"${i}".yaml"
rm -f ${appName}
rm apply -f ${currentServiceName}
done