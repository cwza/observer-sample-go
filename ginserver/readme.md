
## Docs
* https://github.com/gin-gonic/gin
* https://github.com/Depado/ginprom
* https://github.com/opentracing-contrib/go-gin

## Run
``` sh
port="8081" trace_collector_url="http://localhost:14268/api/traces" go run main.go
port="8081" trace_agent_url="localhost:6831" go run main.go
```