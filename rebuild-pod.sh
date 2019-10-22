#!/bin/bash

# This will delete then recompile the Go app, rebuild the container, then deploy to minikube.
# Please verify your environment and dependencies are configured before running this script and deleting the app
# in the process. 


# remove old compiled go app
rm csrgo


# delete any previous csrgo pod
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

