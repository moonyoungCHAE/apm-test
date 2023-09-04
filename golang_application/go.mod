module apm-test-application/golang_application

go 1.20

require (
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gorilla/mux v1.8.0
	github.com/openzipkin-contrib/zipkin-go-sql v0.0.0-20230824044006-7b30542ea014
	github.com/openzipkin/zipkin-go v0.4.2
)

require google.golang.org/grpc v1.57.0 // indirect
