kubectl run -i -t csrgo \
--namespace=jx \
--generator=run-pod/v1 \
--image=csrgo \
--image-pull-policy=Never \
--restart=Never

kubectl logs -n jx csrgo istio-proxy
