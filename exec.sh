kubectl run -i -t csrgo \
   --generator=run-pod/v1 \
   --image=csrgo \
   --image-pull-policy=Never \
   --restart=Never

kubectl logs csrgo
