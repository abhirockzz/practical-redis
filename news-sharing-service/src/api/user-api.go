package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"
)

const usersSet string = "users"

func UserRegistrationHandler(c *gin.Context) {

	userByte, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("registering new user...", string(userByte))

	clientCtxObj, _ := c.Get("redis")
	client := clientCtxObj.(*redis.Client)

	num, err := client.SAdd(usersSet, string(userByte)).Result()
	if err != nil {
		fmt.Println("Error adwhile registering new user", err)
		c.AbortWithError(500, err)
		return
	}

	if num == 0 {
		fmt.Println("user already exists")
		c.Status(http.StatusConflict)
		return
	}

	fmt.Println("registered user successfully")
	c.Status(http.StatusNoContent)
}
