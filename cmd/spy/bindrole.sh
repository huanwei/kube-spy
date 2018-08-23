#!/usr/bin/env bash

template=`cat ./rolebinding.yaml`

printf "nodeName=\"$1\"\n\
\ncat << EOF\n$template\nEOF" | bash > ./tmp.yaml

kubectl apply -f tmp.yaml
rm tmp.yaml