
## Docs
* https://github.com/slok/go-http-metrics
* https://github.com/opentracing/opentracing-go

## Install Jaeger
see ../jaegerserver

## Run
``` sh
kubectl port-forward $(kubectl -n observability get pod | grep "my-jaeger" | awk '{print $1}') 14268:14268 -n observability

port="8080" prom_port="9090" trace_collector_url="http://localhost:14268/api/traces" ginserver_url="http://localhost:8081" go run main.go
# port="8080" prom_port="9090" trace_agent_url="localhost:6831" ginserver_url="http://localhost:8081" go run main.go

curl http://localhost:8080/run
```

## Test
* Check prometheus endpoint is ok `curl http://localhost:9090/metrics`
* Check jaeger interface, run following and goto `http://localhost:16686` by your browser
``` sh
kubectl port-forward $(kubectl -n observability get pod | grep "my-jaeger" | awk '{print $1}') 16686:16686 -n observability
```