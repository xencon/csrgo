#!/bin/bash

# Set minikube env
eval $(minikube docker-env)

# Remove the old binary
rm csrgo

# Delete any previous csrgo pods
kubectl delete pod csrgo --grace-period=0 --force  -n jx

# Build the container
GOOS=linux go build -o ./csrgo .
docker build -t csrgo .

kubectl create clusterrolebinding default-view \
--clusterrole=view \
--serviceaccount=default:default

kubectl apply -f fabric8-rbac.yaml
kubectl apply -f service.yaml

kubectl run -i -t csrgo \
--namespace=jx \
--generator=run-pod/v1 \
--image=csrgo \
--image-pull-policy=Never \
--restart=Never \
--replicas=1
