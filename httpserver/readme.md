
## Docs
* https://github.com/slok/go-http-metrics
* https://github.com/opentracing/opentracing-go

## Run
``` sh
port="8080" promPort="9090" trace_collector_url="http://localhost:14268/api/traces" gin_server_url="http://localhost:8081" go run main.go
port="8080" promPort="9090" trace_agent_url="localhost:6831" ginserver_url="http://localhost:8081" go run main.go
```