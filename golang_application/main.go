package main

import (
	"apm-test-application/golang_application/sql"
	"fmt"
	"github.com/aws/aws-xray-sdk-go/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
	"log"
	"math/rand"
	"net/http"
	"os"
)

func init() {
	// conditionally load plugin
	if os.Getenv("ENVIRONMENT") != "local" {
		ec2.Init()
	}

	xray.Configure(xray.Config{
		ServiceVersion: "1.2.3",
	})

	xray.SetLogger(xraylog.NewDefaultLogger(os.Stderr, xraylog.LogLevelError))
	os.Setenv("AWS_XRAY_TRACING_NAME", "test_service_name")
	//os.Setenv("AWS_XRAY_DAEMON_ADDRESS", "172.17.0.2:2000")
}

func main() {
	if err := sql.MySQLInit(); err != nil {
		log.Println("")
	}

	http.Handle("/", xray.Handler(xray.NewFixedSegmentNamer("MyApp"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!"))
	})))

	http.Handle("/ping", xray.Handler(xray.NewFixedSegmentNamer("MyApp"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PingFunction1("ping start")
		w.Write([]byte("pong!"))
	})))

	http.Handle("/golang", xray.Handler(xray.NewFixedSegmentNamer("MyApp"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("golang application received messag"))
	})))

	http.Handle("/golang/db", xray.Handler(xray.NewFixedSegmentNamer("MyApp"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := sql.GetFruits(r.Context())
		w.Write([]byte(msg))
	})))

	http.ListenAndServe(":8081", nil)
}

func PingFunction1(arg1 string) {
	log.Println(fmt.Sprintf("ping function 1 received (arg1: %s)", arg1))
	PingFunction2(arg1, rand.Int())
}

func PingFunction2(arg1 string, arg2 int) {
	log.Println(fmt.Sprintf("ping function 2 received (arg1: %s, arg2: %d)", arg1, arg2))
	return
}
