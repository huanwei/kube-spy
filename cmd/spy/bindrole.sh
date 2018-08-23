#!/usr/bin/env bash

kubectl create rolebinding spy-view-binding --clusterrole=view --user=system:node:${1} --namespace=default
