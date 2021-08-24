package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
)

const serviceName = "ginserver"

var (
	port               string
	traceAgentUrl      string
	traceCollectortUrl string
)

func readEnvs() {
	port = os.Getenv("port")
	traceAgentUrl = os.Getenv("trace_agent_url")
	traceCollectortUrl = os.Getenv("trace_collector_url")
}

// if you want send traffic to agent then please specify agentUrl, if you want send traffic directly to collector then specify collectorUrl
// if you don't specify either you will get an noop tracer
// ex: agentPort: localhost:6831, collectorUrl: http://localhost:14268/api/traces
func initTracer(serviceName string, agentUrl string, collectorUrl string) (io.Closer, error) {
	var cfg jaegercfg.Configuration
	if agentUrl != "" {
		cfg = jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LocalAgentHostPort: agentUrl,
				LogSpans:           true,
			},
		}
	} else if collectorUrl != "" {
		cfg = jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				CollectorEndpoint: collectorUrl,
				LogSpans:          true,
			},
		}
	} else {
		cfg = jaegercfg.Configuration{}
	}

	closer, err := cfg.InitGlobalTracer(
		serviceName,
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(jaegermetrics.NullFactory),
	)
	if err != nil {
		return nil, err
	}
	return closer, nil
}

func main() {
	readEnvs()

	traceCloser, err := initTracer(serviceName, traceAgentUrl, traceCollectortUrl)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	defer traceCloser.Close()

	r := gin.Default()
	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Namespace(serviceName),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	r.Use(p.Instrument())                                 // add prometheus middleware
	r.Use(ginhttp.Middleware(opentracing.GlobalTracer())) // add open-tracing middleware

	r.GET("/api1/:name", api1)
	r.Run(fmt.Sprintf(":%s", port))
}

func api1(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "api1")
	defer span.Finish()
	span.LogFields(
		opentracinglog.String("event", "soft error"),
		opentracinglog.String("type", "cache timeout"),
		opentracinglog.Int("waited.millis", 1500),
	)
	span.SetTag("tag1", "tag1value")

	name := c.Param("name")
	str := hello(ctx, name)
	c.String(http.StatusOK, str)
}

func hello(ctx context.Context, name string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "hello")
	defer span.Finish()
	return fmt.Sprintf("Hello %s", name)
}
