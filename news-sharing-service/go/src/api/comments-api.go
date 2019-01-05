package api

import (
	"fmt"
	"io/ioutil"
	"model"
	"net/http"

	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"
)

//PostCommentsForNewsItemHandler ...
func PostCommentsForNewsItemHandler(c *gin.Context) {
	commentsBytes, _ := ioutil.ReadAll(c.Request.Body)
	newsID := c.Param("newsid")
	fmt.Println("adding comments for news item " + newsID + " - " + string(commentsBytes))

	clientCtxObj, _ := c.Get("redis")
	client := clientCtxObj.(*redis.Client)

	commentsListName := "news:" + newsID + ":comments"
	_, lpushErr := client.LPush(commentsListName, string(commentsBytes)).Result()

	if lpushErr != nil {
		fmt.Println("Error posting comments for news item", lpushErr)
		c.AbortWithError(500, lpushErr)
		return
	}

	c.Status(http.StatusNoContent)
	fmt.Println("added comment successfully")
}

//GetCommentsForNewsID ...
func GetCommentsForNewsID(c *gin.Context) {
	newsID := c.Param("newsid")
	fmt.Println("getting no. of comments for news ID", newsID)

	commentListName := "news:" + newsID + ":comments"

	clientCtxObj, _ := c.Get("redis")
	client := clientCtxObj.(*redis.Client)

	numOfComments, llenErr := client.LLen(commentListName).Result()
	if llenErr != nil {
		fmt.Println("Error fetching no. of comments", llenErr)
		return
	}

	if numOfComments == 0 {
		fmt.Println("no comments for news ID", newsID)
		c.Data(200, "text/plain", []byte("no comments found"))
		return
	}

	comments, lrangeErr := client.LRange(commentListName, 0, numOfComments).Result()
	if lrangeErr != nil {
		fmt.Println("Error fetching comments", lrangeErr)
		c.AbortWithError(500, lrangeErr)
		return
	}

	c.JSON(200, model.NewsItemComments{newsID, comments})
}
