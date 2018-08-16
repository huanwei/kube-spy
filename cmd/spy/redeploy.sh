#!/usr/bin/env bash
kubectl delete configmap spy-config
kubectl create configmap spy-config --from-file=spy=config.yaml
kubectl delete -f spy-deploy.yaml
kubectl apply -f spy-deploy.yaml