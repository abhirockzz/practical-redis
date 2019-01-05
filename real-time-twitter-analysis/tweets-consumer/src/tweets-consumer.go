package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const tweetRedisListName = "tweets"
const tweetsProcessorListName = "tweets-processing-queue"

//1. reliable transfer (RPOPLPUSH) from list (tweets) to another list (tweets-processing-queue)
//2. processing of tweet i.e. populating tweet details in HASHes and SETs
var client *redis.Client

func main() {

	redisServer := getFromEnvOrDefault("REDIS_HOST", "localhost")
	redisPort := getFromEnvOrDefault("REDIS_PORT", "6379")
	client = redis.NewClient(&redis.Options{Addr: redisServer + ":" + redisPort})
	_, pingErr := client.Ping().Result()
	if pingErr != nil {
		fmt.Println("could not connect to Redis due to " + pingErr.Error())
		return
	}
	defer client.Close()

	for {
		tweetJSON, err := client.BRPopLPush(tweetRedisListName, tweetsProcessorListName, 0*time.Second).Result()
		if err != nil {
			fmt.Println("failed to push tweet info to "+tweetsProcessorListName, err.Error())
		} else {
			go process(tweetJSON) //done in a different goroutine
		}
	}

}

const redisSetNamePrefix = "keyword_tweets:"

func process(tweetJSON string) {
	var tweetObj tweet
	unmarshalErr := json.Unmarshal([]byte(tweetJSON), &tweetObj)
	if unmarshalErr == nil {
		fmt.Println("converted tweet to JSON", tweetObj)
	}
	if len(tweetObj.Terms) == 0 {
		return
	}
	hashName := "tweet:" + tweetObj.TweetID
	pipe := client.Pipeline()
	pipe.HMSet(hashName, tweetObj.toMap())

	for _, term := range tweetObj.Terms {
		set1Name := redisSetNamePrefix + term
		pipe.SAdd(set1Name, tweetObj.TweetID)

		set2Name := redisSetNamePrefix + term + ":" + tweetObj.CreatedDate
		pipe.SAdd(set2Name, tweetObj.TweetID).Result()

	}

	_, pipeErr := pipe.Exec()

	if pipeErr != nil {
		fmt.Println("Pipeline execution error " + pipeErr.Error())
	} else {
		fmt.Println("Stored tweet data for analysis")
		_, lRemErr := client.LRem(tweetsProcessorListName, 0, tweetJSON).Result()
		if lRemErr != nil {
			fmt.Println("unable to delete entry from list " + lRemErr.Error())
		}
	}

}

type tweet struct {
	TweetID     string   `json:"tweetID"`
	Tweeter     string   `json:"tweeter"`
	Tweet       string   `json:"tweet"`
	Terms       []string `json:"terms"`
	CreatedDate string   `json:"createdDate"`
}

func (tweet *tweet) toMap() map[string]interface{} {
	tweetDetailMap := make(map[string]interface{})
	tweetDetailMap["tweet_id"] = tweet.TweetID
	tweetDetailMap["tweeter"] = tweet.Tweeter
	tweetDetailMap["tweet"] = tweet.Tweet
	//iterate over tweet.Terms and convert in comma-separated string i.e. [x,y] = x,y
	termsStr := ""
	for _, term := range tweet.Terms {
		termsStr = termsStr + term + ","
	}
	termsStrFinal := strings.TrimRight(termsStr, ",")
	tweetDetailMap["terms"] = termsStrFinal
	tweetDetailMap["created_date"] = tweet.CreatedDate
	return tweetDetailMap
}
func getFromEnvOrDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}

	return val
}
