#!/usr/bin/env bash
kubectl delete rolebinding spy-deploy-binding
kubectl create rolebinding spy-deploy-binding --clusterrole=system:controller:deployment-controller  --user=system:node:${1} --namespace=default
kubectl delete rolebinding spy-pod-binding
kubectl create rolebinding spy-pod-binding --clusterrole=edit  --user=system:node:${1} --namespace=default