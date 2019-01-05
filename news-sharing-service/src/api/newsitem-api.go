package api

import (
	"errors"
	"fmt"
	"model"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"
)

const redisNewsIdCounter string = "news-id-counter"
const redisNewsHashPrefix string = "news"
const redisNewsUpvotesSortedSet string = "news-upvotes"

//PostNewsItemHandler ...
func PostNewsItemHandler(c *gin.Context) {

	var news model.NewItemSubmission
	c.ShouldBindJSON(&news)

	submittedBy := c.GetHeader("user")

	fmt.Println("Trying to add news item\n", news)

	clientCtxObj, _ := c.Get("redis")
	client := clientCtxObj.(*redis.Client)

	exists, sisMemberErr := client.SIsMember(usersSet, submittedBy).Result()

	if sisMemberErr != nil {
		fmt.Println("Error trying to check user membership", sisMemberErr)
		c.AbortWithError(500, sisMemberErr)
		return
	}

	if !exists {
		fmt.Println("Invalid user", submittedBy)
		c.AbortWithError(500, errors.New("Invalid user "+submittedBy))
		return
	}

	newsID, err := client.Incr(redisNewsIdCounter).Result()
	if err != nil {
		fmt.Println("Error trying to generate new item ID", err)
		c.AbortWithError(500, err)
		return
	}

	fmt.Println("generated news ID", newsID)

	newsHashKey := redisNewsHashPrefix + ":" + strconv.Itoa(int(newsID))
	fmt.Println("Saving news details in", newsHashKey)

	newsDetails := map[string]interface{}{"title": news.Title, "url": news.Url, "submittedBy": submittedBy}
	_, hmsetErr := client.HMSet(newsHashKey, newsDetails).Result()

	if hmsetErr != nil {
		fmt.Println("Error saving news item", hmsetErr)
		c.AbortWithError(500, hmsetErr)
		return
	}

	_, zaddErr := client.ZAdd(redisNewsUpvotesSortedSet, redis.Z{Member: newsID, Score: 0}).Result()

	if zaddErr != nil {
		fmt.Println("Error adding upvotes for news ID", zaddErr)
		c.AbortWithError(500, zaddErr)
		return
	}

	fmt.Println("saved news item successfully", newsHashKey)
	c.JSON(http.StatusCreated, newsID)
}

//GetNewsItemByID ...
func GetNewsItemByIDHandler(c *gin.Context) {
	newsID := c.Param("newsid")
	fmt.Println("searching for news item with ID", newsID)

	clientCtxObj, _ := c.Get("redis")
	client := clientCtxObj.(*redis.Client)

	newsItem, err := getNewsItemDetails(newsID, client)
	if err != nil {
		fmt.Println("Could not find news item", err)
		c.AbortWithError(500, err)
		return
	}

	if newsItem.NewsID == "" {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, newsItem)
	}

}

//GetAllNewsItemsHandler ...
func GetAllNewsItemsHandler(c *gin.Context) {
	fmt.Println("listing all news items...")

	clientCtxObj, _ := c.Get("redis")
	client := clientCtxObj.(*redis.Client)

	numItems, zcardErr := client.ZCard(redisNewsUpvotesSortedSet).Result()

	if zcardErr != nil {
		fmt.Println("Error finding no. of entries in new item upvotes sorted set", zcardErr)
		c.AbortWithError(500, zcardErr)
		return
	}

	fmt.Println("no. of news items", numItems)
	itemIDs, zrevrangeErr := client.ZRevRange(redisNewsUpvotesSortedSet, 0, numItems-1).Result()
	fmt.Println("no. of news items --", len(itemIDs))

	if zrevrangeErr != nil {
		fmt.Println("Error sorting news item upvoted sorted set", zrevrangeErr)
		c.AbortWithError(500, zrevrangeErr)
		return
	}

	newsItems := []model.NewsItem{} //init to empty slice
	for _, newsID := range itemIDs {
		newsItem, err := getNewsItemDetails(newsID, client)
		if err == nil {
			newsItems = append(newsItems, newsItem)
		}
	}

	fmt.Println("no. of news items found", len(newsItems))
	c.JSON(http.StatusOK, newsItems)
}

func getNewsItemDetails(newsID string, client *redis.Client) (model.NewsItem, error) {
	newsItemHashKey := redisNewsHashPrefix + ":" + newsID
	fmt.Println("searching for news details in HASH", newsItemHashKey)

	var item model.NewsItem
	newsItemDetailMap, hgetallErr := client.HGetAll(newsItemHashKey).Result()

	if hgetallErr != nil {
		fmt.Println("Error trying to find news item details", hgetallErr)
		return item, hgetallErr
	}

	if len(newsItemDetailMap) == 0 {
		fmt.Println("No news item with ID", newsID)
		return item, nil
	}

	upvotes, zscoreErr := client.ZScore(redisNewsUpvotesSortedSet, newsID).Result()

	if zscoreErr != nil {
		fmt.Println("Error finding upvotes for news item", zscoreErr)
		return item, zscoreErr
	}

	commentListName := "news:" + newsID + ":comments"
	fmt.Println("Getting comments for news", newsID)

	comments, llenErr := client.LLen(commentListName).Result()
	if llenErr != nil {
		fmt.Println("Error finding no. of comments for news item", llenErr)
		return item, llenErr
	}

	item = model.NewsItem{newsID, newsItemDetailMap["title"], newsItemDetailMap["submittedBy"], newsItemDetailMap["url"], strconv.Itoa(int(upvotes)), strconv.Itoa(int(comments))}
	fmt.Println("news item details", item)
	return item, nil
}
