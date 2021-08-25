
## Install jaeger, prometheus
see ./jaegerserver and ./promserver

## Install observer-sample-go
``` sh
kubectl create namespace try
cd helm
helm install -n try -f values.yaml observer-sample-go .
# helm delete -n try observer-sample-go
```

## Run
``` sh
curl http://[httpserver nodeip]:[httpserver nodeport]/run
```

## Test
``` sh
kubectl port-forward $(kubectl -n observability get pod | grep "my-jaeger" | awk '{print $1}') 16686:16686 -n observability
kubectl port-forward prometheus-prometheus-kube-prometheus-prometheus-0 9090:9090 -n prometheus
```
* go to `http://localhost:16686` and `http://localhost:9090` by browser

