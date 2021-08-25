
## Docs
* https://github.com/gin-gonic/gin
* https://github.com/Depado/ginprom
* https://github.com/opentracing-contrib/go-gin

## Install Jaeger Server
see ../jaegerserver

## Run
``` sh
kubectl port-forward $(kubectl -n observability get pod | grep "my-jaeger" | awk '{print $1}') 14268:14268 -n observability

port="8081" trace_collector_url="http://localhost:14268/api/traces" go run main.go
# port="8081" trace_agent_url="localhost:6831" go run main.go

curl http://localhost:8081/api1/xx
```

## Test
* Check prometheus endpoint is ok `curl http://localhost:8081/metrics`
* Check jaeger interface, run following and goto `http://localhost:16686` by your browser
``` sh
kubectl port-forward $(kubectl -n observability get pod | grep "my-jaeger" | awk '{print $1}') 16686:16686 -n observability
```