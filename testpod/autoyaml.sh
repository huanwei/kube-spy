#!/usr/bin/env bash
template1=`cat ./http-test-template.yaml`
template2=`cat ./http-test-service-template.yaml`

for((i=1;i<=${1};i++));
do
appName="http-test"${i}
currentServiceName="http-test-service"${i}
if [ $i -eq ${1} ]
then
    nextServiceName=""
else
    nextServiceName="http-test-service"$(($i+1))
fi
printf "appName=\"${appName}\"\n\
nextServiceName=\"${nextServiceName}\"\n\
currentServiceName=\"${currentServiceName}\"\n\
\ncat << EOF\n$template1\nEOF" | bash > ./${appName}.yaml

printf "currentServiceName=\"${currentServiceName}\"\n\
appName=\"${appName}\"\n\
\ncat << EOF\n$template2\nEOF" | bash > ./${currentServiceName}.yaml

done

