#!/usr/bin/env bash
kubectl delete configmap spy-config -n default
kubectl create configmap spy-config --from-file=spy=config.yaml -n default
kubectl delete -f spy-job.yaml
kubectl apply -f spy-job.yaml