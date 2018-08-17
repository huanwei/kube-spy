#!/usr/bin/env bash
echo " ">/tmp/watch
time=0.0
for ((i=1;i<=30;i++));do
echo "$time s:">>/tmp/watch
kubectl describe pod http|grep chaos= >>/tmp/watch
sleep 0.2
time=$(echo "$time + 0.2"|bc)
done


