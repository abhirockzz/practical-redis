package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"github.com/dghubble/oauth1"

	"github.com/dghubble/go-twitter/twitter"
)

const redisTweetListName = "tweets"

var apiStopChannel chan interface{}
var active bool

type tweetInfo struct {
	TweetID     string   `json:"tweetID"`
	Tweeter     string   `json:"tweeter"`
	Tweet       string   `json:"tweet"`
	Terms       []string `json:"terms"`
	CreatedDate string   `json:"createdDate"`
}

//Start ...
func Start() error {

	redisServer := getFromEnvOrDefault("REDIS_HOST", "localhost")
	redisPort := getFromEnvOrDefault("REDIS_PORT", "6379")
	redisClient := redis.NewClient(&redis.Options{Addr: redisServer + ":" + redisPort})
	_, pingErr := redisClient.Ping().Result()
	if pingErr != nil {
		fmt.Println("could not connect to Redis due to " + pingErr.Error())
		return pingErr
	}

	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return errors.New("Please specify valid Twitter credentials")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	twitterClient := twitter.NewClient(httpClient)

	fmt.Println("Twitter Client setup")

	trackedTerms := os.Getenv("TWITTER_TRACKED_TERMS")
	//trackedTerms := getFromEnvOrDefault("TWITTER_TRACKED_TERMS", "trump,realDonaldTrump,java,go,redis")

	trackedTermsSlice := strings.Split(trackedTerms, ",")
	params := &twitter.StreamFilterParams{
		Track:         trackedTermsSlice,
		StallWarnings: twitter.Bool(true),
	}
	var err error
	stream, err := twitterClient.Streams.Filter(params)
	if err != nil {
		return err
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
		matches := getMatchedTerms(tweet.Text)
		date := formatTweetDate(tweet.CreatedAt)
		tweetInfoByte, marshalErr := json.Marshal(tweetInfo{TweetID: strconv.Itoa(int(tweet.ID)), Tweeter: tweet.User.ScreenName, Tweet: tweet.Text, Terms: matches, CreatedDate: date})
		if marshalErr != nil {
			fmt.Println("failed to marshal TweetInfo to JSON", marshalErr)
		} else {
			go func() {
				_, lpushErr := redisClient.LPush(redisTweetListName, string(tweetInfoByte)).Result()
				if lpushErr != nil {
					fmt.Println("failed to push tweet info to Redis", lpushErr)
				}
			}()
		}

		time.Sleep(10 * time.Second) //TODO remove
	}

	go demux.HandleChan(stream.Messages)
	active = true
	fmt.Println("Started listening to twitter stream for terms - ", trackedTerms)

	apiStopChannel = make(chan interface{})

	go func() {
		fmt.Println("Waiting for listener to stop")
		<-apiStopChannel
		fmt.Println("Listener stop request")
		active = false

		stream.Stop()
		fmt.Println("Listener stopped...")

		redisClient.Close()
		fmt.Println("Redis connection closed...")
	}()

	return nil

}

func GetTweetsListenerStatus() bool {
	return active
}

func Stop() {
	apiStopChannel <- "stop"
}

func getMatchedTerms(tweet string) []string {
	var matches []string
	wordsInTweet := strings.Split(tweet, " ")
	trackedTerms := getFromEnvOrDefault("TWITTER_TRACKED_TERMS", "trump,realDonaldTrump,potus,java,golang,redis")
	trackedTermsSlice := strings.Split(trackedTerms, ",")
	for _, trackedTerm := range trackedTermsSlice {
		for _, word := range wordsInTweet {
			wordL := strings.ToLower(word)
			trackedTermL := strings.ToLower(trackedTerm)
			if wordL == trackedTermL || wordL == "@"+trackedTermL || wordL == "#"+trackedTermL {
				matches = append(matches, trackedTerm)
			}
		}
	}

	//remove duplicates - https://gist.github.com/alioygur/16c66b4249cb42715091fe010eec7e33
	for i := 0; i < len(matches); i++ {
		for i2 := i + 1; i2 < len(matches); i2++ {
			if matches[i] == matches[i2] {
				// delete
				matches = append(matches[:i2], matches[i2+1:]...)
				i2--
			}
		}
	}

	return matches
}

func getFromEnvOrDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}

	return val
}

func formatTweetDate(date string) string {
	pTime, _ := time.Parse(time.RubyDate, date)
	return pTime.Format("02-01-2006")
}
