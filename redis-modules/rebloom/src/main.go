package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reco"
	"strconv"
	"util"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var redisCoordinate string

func init() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" {
		redisHost = "192.168.99.100"
	}

	if redisPort == "" {
		redisPort = "6379"
	}

	redisCoordinate = redisHost + ":" + redisPort
	fmt.Println("Redis server - " + redisCoordinate)
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/load/", loadData).Methods("GET")
	r.HandleFunc("/{user}/article/", readArticle).Methods("GET")
	r.HandleFunc("/{user}/articles/", getRecommendedArticles).Methods("GET")

	appPort := os.Getenv("PORT")

	if appPort == "" {
		appPort = "8080"
	}

	fmt.Println("starting recommendation service on port " + appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, r))
}

func loadData(resp http.ResponseWriter, req *http.Request) {
	util.LoadArticlesForTopics(redisCoordinate)
	util.LoadUserInterests(redisCoordinate)
	resp.Write([]byte("Added articles to topics and user interests. Execute KEYS * using redis-cli to check"))
}

const redisBloomFilterNamePrefix = "RecommendationHits-"
const redisBloomFilterAddCmd = "BF.ADD"

func readArticle(resp http.ResponseWriter, req *http.Request) {
	article := req.URL.Query()["url"][0]
	user := mux.Vars(req)["user"]
	fmt.Println("getting article " + article + " for user " + user)
	//get the article
	res, err := http.Get(article)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	contents, err := ioutil.ReadAll(res.Body)

	res.Body.Close()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	contentType := res.Header.Get("Content-Type")
	resp.Header().Set("Content-Type", contentType)

	_, writeErr := resp.Write(contents)

	if writeErr != nil {
		http.Error(resp, writeErr.Error(), http.StatusInternalServerError)
		return
	}

	//if article has been accessed, we can it add to bloom filter (specific to the user)
	client := redis.NewClient(&redis.Options{Addr: redisCoordinate})
	_, perr := client.Ping().Result()

	if perr != nil {
		fmt.Println("Could not connect to Redis " + perr.Error())
	} else {
		fmt.Println("Connected to Redis " + redisCoordinate)
	}
	defer client.Close()

	redisBloomFilterName := redisBloomFilterNamePrefix + user
	fmt.Println("Adding read article to user's bloom filter " + redisBloomFilterName)

	addToBloomFilterCmd := redis.NewIntCmd(redisBloomFilterAddCmd, redisBloomFilterName, article)
	client.Process(addToBloomFilterCmd)
	cmdResult, _ := addToBloomFilterCmd.Result()
	fmt.Println("Article marked as read " + strconv.Itoa(int(cmdResult)))

	fmt.Println("closed redis connection...")
}

func getRecommendedArticles(resp http.ResponseWriter, req *http.Request) {
	user := mux.Vars(req)["user"]
	fmt.Println("User -" + user)
	recoUtil := reco.NewRecommendationUtil(redisCoordinate)
	defer recoUtil.CloseConn()

	recommendedArticles := recoUtil.GenArticleRecommendations(user)
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(recommendedArticles)
}
