package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

//RedisConnHandler - middleware to establish and close Redis connection
func RedisConnHandler(c *gin.Context) {

	fmt.Println("connecting to Redis...")
	connected := true

	redisServer := os.Getenv("REDIS_SERVER")
	if redisServer == "" {
		redisServer = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{Addr: redisServer})
	err := client.Ping().Err()
	if err != nil {
		fmt.Println("Could not connect to Redis...", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		connected = false
	}

	fmt.Println("Connected to Redis server at", redisServer)

	c.Set("redis", client)

	c.Next()

	if client != nil && connected {
		client.Close()
		fmt.Println("Closed Redis connection...")
	}
}
