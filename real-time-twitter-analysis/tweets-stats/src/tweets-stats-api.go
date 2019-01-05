package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func connectToRedis() (*redis.Client, error) {

	redisServer := getFromEnvOrDefault("REDIS_HOST", "localhost")
	redisPort := getFromEnvOrDefault("REDIS_PORT", "6379")
	client := redis.NewClient(&redis.Options{Addr: redisServer + ":" + redisPort})
	_, pingErr := client.Ping().Result()
	if pingErr != nil {
		fmt.Println("could not connect to Redis due to " + pingErr.Error())
		return nil, pingErr
	}

	return client, nil
}

func main() {
	router := gin.Default()
	router.GET("tweets", findTweetsHandler)
	router.Run()
}

const setNamePrefix string = "keyword_tweets:"
const hashNamePrefix string = "tweet:"

//http://localhost:8081/tweets/?keywords=a,b,c&op=OR&date=20-06-2018
func findTweetsHandler(c *gin.Context) {
	fmt.Println("request URL", c.Request.URL.String())

	keywords := c.Query("keywords")

	if keywords == "" {
		c.Status(400)
		c.Writer.WriteString("keywords query parameter cannot be empty")
		return
	}

	fmt.Println("searching for tweets with keywords", keywords)

	date := c.Query("date")

	if date != "" {
		fmt.Println("searching for tweets on", date)
	}
	operation := c.Query("op")

	if operation != "" {
		fmt.Println("applying operation", operation)
	}

	tweets, err := findTweets(keywords, operation, date)

	if err != nil {
		c.Status(500)
		c.Writer.WriteString("Unable to fetch tweets due to " + err.Error())
		return
	}

	c.JSON(200, tweets)
}

func findTweets(keywords, operation, date string) ([]map[string]string, error) {

	var sets []string
	keywordsSlice := strings.Split(keywords, ",")
	for _, keyword := range keywordsSlice {
		var setName string
		if date == "" {
			setName = setNamePrefix + keyword
		} else {
			setName = setNamePrefix + keyword + ":" + date
		}
		sets = append(sets, setName)
	}

	fmt.Println("Keyword sets", sets)
	var tempSetName string
	deleteSet := true
	client, connErr := connectToRedis()
	if connErr != nil {
		fmt.Println("Could not connect to Redis", connErr.Error())
		return nil, connErr
	}
	defer client.Close()
	switch operation {
	case "":
		fmt.Println("no operation specified")
		tempSetName = sets[0]
		deleteSet = false
	case "AND":
		fmt.Println("AND operation specified")
		tempSetName, _ = generateRandomString(10)
		client.SInterStore(tempSetName, sets...)
	case "OR":
		fmt.Println("OR operation specified")
		tempSetName = sets[0]
		tempSetName, _ = generateRandomString(10)
		client.SUnionStore(tempSetName, sets...)
	}

	fmt.Println("Temporary set name", tempSetName)

	tweetIDs, smembersErr := client.SMembers(tempSetName).Result()
	if smembersErr != nil {
		return nil, errors.New("Unable to find members of set " + tempSetName)
	}

	var tweets []map[string]string
	for _, tweetID := range tweetIDs {
		tweetInfoHashName := hashNamePrefix + tweetID
		tweetInfoMap, hgetallErr := client.HGetAll(tweetInfoHashName).Result()
		if hgetallErr != nil {
			fmt.Println("unable to fetch info for tweet ID", tweetID)
		}
		tweets = append(tweets, tweetInfoMap)
	}

	if tempSetName != "" && deleteSet {
		_, delErr := client.Del(tempSetName).Result()

		if delErr != nil {
			fmt.Println("unable to delete set", tempSetName)
		} else {
			fmt.Println("deleted set", tempSetName)

		}
	}
	fmt.Println("No. of  tweets found", len(tweets))
	return tweets, nil

}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func getFromEnvOrDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}

	return val
}
