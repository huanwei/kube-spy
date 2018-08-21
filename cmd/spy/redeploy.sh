#!/usr/bin/env bash
kubectl delete configmap spy-config -n kube-system
kubectl create configmap spy-config --from-file=spy=config.yaml -n kube-system
kubectl delete -f spy-deploy.yaml
kubectl apply -f spy-deploy.yaml