package main

import (
	"apm-test-application/golang_application/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/lambda"
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

	http.Handle("/lambda", xray.Handler(xray.NewFixedSegmentNamer("MyApp"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		var endpoint = "lambda.ap-northeast-2.api.aws"
		// Lambda 클라이언트 생성
		svc := lambda.New(sess, &aws.Config{
			Region:   aws.String("ap-northeast-2"), // AWS Lambda 함수가 있는 리전
			Endpoint: &endpoint,
		})

		// Lambda 함수 호출
		result, err := svc.Invoke(&lambda.InvokeInput{
			FunctionName:   aws.String("tommoy-test"),     // Lambda 함수 이름
			InvocationType: aws.String("RequestResponse"), // 동기식 호출
		})

		if err != nil {
			fmt.Println("Lambda 함수 호출 에러:", err)
			w.Write([]byte(err.Error()))
			return
		}

		// Lambda 함수의 응답 데이터 출력
		fmt.Println("Lambda 함수 응답:", string(result.Payload))
		w.Write([]byte(result.Payload))
	})))

	http.Handle("/dynamo_db", xray.Handler(xray.NewFixedSegmentNamer("MyApp"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("ap-northeast-2"),
		})
		if err != nil {
			fmt.Println("session fail: " + err.Error())
			w.Write([]byte("session fail " + err.Error()))
			return
		}
		svc := dynamodb.New(sess)
		result, err := svc.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String("tommoy-test"),
			Key: map[string]*dynamodb.AttributeValue{
				"title": {
					S: aws.String("title"),
				},
			},
		})
		if err != nil {
			fmt.Println("get item fail: " + err.Error())
			w.Write([]byte("get item fail: " + err.Error()))
			return
		}
		fmt.Println(result)

		var m map[string]interface{}
		if err := dynamodbattribute.UnmarshalMap(result.Item, &m); err != nil {
			fmt.Println("unmarshal error ", err.Error())
			w.Write([]byte("unmarshal error " + err.Error()))
			return
		}
		bytes, err := json.Marshal(&m)
		if err != nil {
			fmt.Println("marshal error ", err.Error())
			w.Write([]byte("marshal error " + err.Error()))
			return
		}
		w.Write(bytes)
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

var testData = `
{  
   "title" : "title",
   "artist" : "yundream",
   "createat" : "2020.01.05",
   "tracklist" : [
      {
         "number" : 1,
         "playtime" : 185,
         "Name" : "Hello world"
      },
      {
         "number" : 2,
         "playtime" : 212,
         "Name" : "Wonderful world"
      },
      {
         "number" : 3,
         "Name" : "Hell of Fire",
         "playtime" : 198
      }
   ]
}
`
