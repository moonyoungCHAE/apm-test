package main

import (
	"apm-test-application/golang_application/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	sql.MySQLInit()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		PingFunction1("ping start")
		c.String(http.StatusOK, "pong")
	})

	r.GET("/golang", func(c *gin.Context) {
		c.JSON(http.StatusOK, "golang application received message")
	})

	r.GET("/golang/db", func(c *gin.Context) {
		msg := sql.GetFruits()
		c.JSON(http.StatusOK, msg)
	})

	r.Run(":8081")
}

func PingFunction1(arg1 string) {
	log.Println(fmt.Sprintf("ping function 1 received (arg1: %s)", arg1))
	PingFunction2(arg1, rand.Int())
}

func PingFunction2(arg1 string, arg2 int) {
	log.Println(fmt.Sprintf("ping function 2 received (arg1: %s, arg2: %d)", arg1, arg2))
	return
}
