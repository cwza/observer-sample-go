
``` sh
kubectl port-forward $(kubectl -n observability get pod | grep "my-jaeger" | awk '{print $1}') 14268:14268 -n observability
kubectl port-forward $(kubectl -n observability get pod | grep "my-jaeger" | awk '{print $1}') 16686:16686 -n observability
```

``` sh
kubectl create namespace try
cd helm
helm install -n try -f values.yaml observer-sample-go .
helm delete -n try observer-sample-go
```