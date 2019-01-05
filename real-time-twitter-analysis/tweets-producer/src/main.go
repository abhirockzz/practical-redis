package main

import (
	"lcm"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("tweets/producer", lcm.StartServiceHandler)
	router.DELETE("tweets/producer", lcm.StopServiceHandler)

	router.Run()
}
