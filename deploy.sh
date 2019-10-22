#!/bin/bash

# Set minikube env
eval $(minikube docker-env)

# Remove the old binary
rm csrgo

# Create namespace
#kubectl create ns csrgo

# Delete any previous csrgo pods
kubectl delete pod csrgo --grace-period=0 --force  -n default
	
# Build the container
GOOS=linux go build -o ./app .
docker build -t csrgo .

kubectl create clusterrolebinding default-view \
--clusterrole=view \
--serviceaccount=default:default

kubectl apply -f csrgo-rbac.yaml
#kubectl apply -f service.yaml

kubectl run -i -t csrgo \
--namespace=default \
--generator=run-pod/v1 \
--image=csrgo \
--image-pull-policy=Never \
--restart=Never \
--replicas=1

kubectl -n default logs csrgo
