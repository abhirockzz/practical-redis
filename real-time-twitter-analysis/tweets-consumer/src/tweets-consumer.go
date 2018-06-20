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

var client *redis.Client

func init() {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" {
		redisHost = "192.168.99.100"
	}

	if redisPort == "" {
		redisPort = "6379"
	}

	redisCoordinate := redisHost + ":" + redisPort
	fmt.Println("Redis server - " + redisCoordinate)

	client = redis.NewClient(&redis.Options{Addr: redisCoordinate})
	_, perr := client.Ping().Result()

	if perr != nil {
		fmt.Println("Could not connect to Redis " + perr.Error())
	} else {
		fmt.Println("Consumer connected to Redis...")
	}
}

//1. reliable transfer (RPOPLPUSH) from list (tweets) to another list (tweets-processing-queue)
//2. processing of tweet i.e. populating tweet details in HASHes and SETs

func main() {

	defer client.Close()

	for {
		tweetJSON, err := client.BRPopLPush(tweetRedisListName, tweetsProcessorListName, 0*time.Second).Result()
		if err != nil {
			fmt.Println("failed to push tweet info to " + tweetsProcessorListName)
		} else {
			//fmt.Println("pushed " + tweetJSON + " to " + tweetsProcessorListName)
		}

		go process(tweetJSON) //done in a different goroutine
	}

}

const redisSetNamePrefix = "keyword_tweets:"

func process(tweetJSON string) {
	var tweetObj tweet
	unmarshalErr := json.Unmarshal([]byte(tweetJSON), &tweetObj)
	if unmarshalErr == nil {
		fmt.Println("converted tweet to JSON", tweetObj)
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
		//fmt.Println("Pipeline executed successfully")
		_, lRemErr := client.LRem(tweetsProcessorListName, 0, tweetJSON).Result()
		if lRemErr != nil {
			fmt.Println("unable to delete entry from list " + lRemErr.Error())
		} else {
			//fmt.Println("no. of elements removed ", num)
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
	tweetDetailMap["text"] = tweet.Tweet
	//iterate over tweet.Terms and convert in comma-separated string i.e. [x,y] = x,y
	termsStr := ""
	for _, term := range tweet.Terms {
		termsStr = termsStr + term + ","
	}
	termsStrFinal := strings.TrimRight(termsStr, ",")
	//fmt.Println("terms -- " + termsStrFinal)
	tweetDetailMap["terms"] = termsStrFinal
	tweetDetailMap["created_date"] = tweet.CreatedDate
	return tweetDetailMap
}
