# Goreport

## Client-go based In-cluster container level resource reporting

##### Clear any previous minikube cluster and 
##### delete ~/.minikube ~/.jx ~/.kube ~/.helm
##### and start a new minikube cluster
```
minikube start
```
##### Build/Rebuild client-go app 
```
./rebuild-pod.sh
```
##### Run the pod again without rebuilding
```
./run-pod.sh
```
