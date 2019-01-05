package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func UpvoteNewsItemHandler(c *gin.Context) {
	newsID := c.Param("newsid")
	fmt.Println("Upvoting news ID", newsID)

	clientCtxObj, _ := c.Get("redis")
	client := clientCtxObj.(*redis.Client)

	_, zincrErr := client.ZIncr(redisNewsUpvotesSortedSet, redis.Z{Member: newsID, Score: 1}).Result()

	if zincrErr != nil {
		fmt.Println("Error while upvoting news item", zincrErr)
		c.AbortWithError(500, zincrErr)
		return
	}

	fmt.Println("upvoted news item successfully")
	c.Status(http.StatusNoContent)
}

func GetUpvotesForNewsItemHandler(c *gin.Context) {
	newsID := c.Param("newsid")
	fmt.Println("geting no. of upvotes for news ID", newsID)

	clientCtxObj, _ := c.Get("redis")
	client := clientCtxObj.(*redis.Client)

	upvotes, zscoreErr := client.ZScore(redisNewsUpvotesSortedSet, newsID).Result()

	if zscoreErr != nil {
		fmt.Println("Error fetching no. of upvotes for news item", zscoreErr)
		c.AbortWithError(500, zscoreErr)
		return
	}

	fmt.Println("no. of upvotes", upvotes)
	c.Writer.WriteString(strconv.Itoa((int(upvotes))))

}
