#!/bin/bash

# Remove the old binary
rm csrgo


# Delete any previous csrgo pods
kubectl delete pod csrgo --grace-period=5  -n default


GOOS=linux go build -o ./csrgo .
eval $(minikube docker-env)
docker build -t csrgo .

kubectl create clusterrolebinding default-view \
   --clusterrole=view \
   --serviceaccount=default:default

kubectl apply -f fabric8-rbac.yaml

kubectl run -i -t csrgo \
   --generator=run-pod/v1 \
   --image=csrgo \
   --image-pull-policy=Never \
   --restart=Never \
   --replicas=1

# kubectl logs csrgo

