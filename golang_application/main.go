package main

import (
	"apm-test-application/golang_application/sql"
	"fmt"
	"github.com/gorilla/mux"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
	"math/rand"
	"net/http"
)

const EndpointURL = "https://tracing-kakao-collector.onkakao.net/api/v2/spans"
const ServiceName = "golang<eeabcf5f48714ef3b75b51e43a252b2e>"

//const EndpointURL = "http://localhost:9411/api/v2/spans" // local zipkin 연동할 때

func main() {
	tracer, err := newTracer()
	if err != nil {
		log.Fatal(err)
	}

	if err := sql.MySQLInit(tracer); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(zipkinhttp.NewServerMiddleware(
		tracer,
		zipkinhttp.SpanName("request")), // name for request span
	)
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello!"))
	})
	r.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		PingFunction1("ping start")
		writer.Write([]byte("pong!"))
	})
	r.HandleFunc("/golang", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("golang application received messag"))
	})
	r.HandleFunc("/golang/db", func(writer http.ResponseWriter, request *http.Request) {
		msg := sql.GetFruits(tracer)
		writer.Write([]byte(msg))
	})

	log.Fatal(http.ListenAndServe(":8081", r))
}

func PingFunction1(arg1 string) {
	log.Println(fmt.Sprintf("ping function 1 received (arg1: %s)", arg1))
	PingFunction2(arg1, rand.Int())
}

func PingFunction2(arg1 string, arg2 int) {
	log.Println(fmt.Sprintf("ping function 2 received (arg1: %s, arg2: %d)", arg1, arg2))
	return
}

func newTracer() (*zipkin.Tracer, error) {
	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(EndpointURL)

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: ServiceName, Port: 8081}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	return t, err
}
