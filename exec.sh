kubectl run -i -t csrgo \
--namespace=csrgo \
--generator=run-pod/v1 \
--image=csrgo \
--image-pull-policy=Never \
--restart=Never

kubectl logs -n csrgo csrgo istio-proxy
