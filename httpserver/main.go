package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	prommetrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prommiddleware "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
)

const serviceName = "httpserver"

var (
	port               string
	promPort           string
	traceAgentUrl      string
	traceCollectortUrl string
	ginServerUrl       string
	grpcServerUrl      string
)

func readEnvs() {
	getNoneEmptyEnv := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			log.Fatalf("failed to getEnv, key: %s", key)
		}
		return val
	}
	port = getNoneEmptyEnv("port")
	promPort = getNoneEmptyEnv("prom_port")
	traceAgentUrl = os.Getenv("trace_agent_url")
	traceCollectortUrl = os.Getenv("trace_collector_url")
	ginServerUrl = getNoneEmptyEnv("ginserver_url")
	grpcServerUrl = getNoneEmptyEnv("grpcserver_url")
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

	mdlw := prommiddleware.New(prommiddleware.Config{
		Recorder: prommetrics.NewRecorder(prommetrics.Config{}),
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/run", run)
	h := std.Handler("", mdlw, mux) // add prometheus middleware

	// Serve metrics
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", promPort), promhttp.Handler()); err != nil {
			log.Panicf("error while serving metrics: %s", err)
		}
	}()

	// Serve our handler.
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), h); err != nil {
		log.Fatalf("error while serving: %s", err)
	}

}

func run(w http.ResponseWriter, r *http.Request) {
	// extract tracer context from request
	wiredCtx, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)
	var span opentracing.Span
	if err != nil {
		span = opentracing.StartSpan("run")
	} else {
		span = opentracing.StartSpan(
			"run",
			ext.RPCServerOption(wiredCtx),
		)
	}

	defer span.Finish()

	// call ginserver
	url := ginServerUrl + "/api1/kevin"
	httpClient := &http.Client{}
	httpReq, _ := http.NewRequest("GET", url, nil)
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(httpReq.Header),
	)
	_, err = httpClient.Do(httpReq)
	if err != nil {
		fmt.Fprintf(w, "err2: %s", err)
	}

	// call grpcserver

	fmt.Fprintln(w, "Success")
}
