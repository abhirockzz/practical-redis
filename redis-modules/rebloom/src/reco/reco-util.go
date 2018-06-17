package reco

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

const redisBloomFilterPrefix = "RecommendationHits:"
const redisBloomFilterExistsCmd = "BF.EXISTS"

//RecommendationUtil ...
type RecommendationUtil struct {
	redisClient redis.Client
}

//NewRecommendationUtil ...
func NewRecommendationUtil(redisCoordinate string) RecommendationUtil {
	client := redis.NewClient(&redis.Options{Addr: redisCoordinate})

	_, err := client.Ping().Result()

	if err != nil {
		fmt.Println("Could not connect to Redis " + err.Error())
	} else {
		fmt.Println("NewRecommendationUtil connected to redis")
	}

	return RecommendationUtil{redisClient: *client}
}

//GenArticleRecommendations - based on SUNION and Bloom Filter to avoid recommending already read articles
func (recoUtil *RecommendationUtil) GenArticleRecommendations(user string) []string {
	//if bloom filter contains a recommnded article, do not include in final set of recommendations
	rawRecos := recoUtil.genRawArticleRecommendations(user)

	fmt.Println("generating fine tuned articles recommendations for user " + user)
	fmt.Println("checking in Bloom Filter for user " + user)

	var finalRecos []string
	for i := 0; i < len(rawRecos); i++ {
		if recoUtil.isArticleAlreadyReadByUser(user, rawRecos[i]) == 0 { //has NOT been read for SURE
			finalRecos = append(finalRecos, rawRecos[i])
		} else {
			fmt.Println("article " + rawRecos[i] + " has already been read")
		}
	}

	fmt.Println("no. of fine tuned recommendations generated = " + strconv.Itoa(len(finalRecos)))
	return finalRecos
}

//this is just based on SUNION (without bloom filter)
func (recoUtil *RecommendationUtil) genRawArticleRecommendations(user string) []string {
	fmt.Println("generating RAW articles recos for user " + user)
	userInteresetsSet := "user:" + user + ":interests"
	fmt.Println("getting interests from set " + userInteresetsSet)

	members, _ := recoUtil.redisClient.SMembers(userInteresetsSet).Result()
	fmt.Println("got " + strconv.Itoa(len(members)) + " interests for user " + user)

	var recoSetArr []string
	for i := 0; i < len(members); i++ {
		topicSetName := "topic:" + members[i] + ":articles"
		recoSetArr = append(recoSetArr, topicSetName)
	}

	recos, _ := recoUtil.redisClient.SUnion(recoSetArr...).Result()
	fmt.Println("no. of raw recommendations generated = " + strconv.Itoa(len(recos)))

	return recos
}

func (recoUtil *RecommendationUtil) isArticleAlreadyReadByUser(user string, article string) int {
	redisBloomFilterName := redisBloomFilterPrefix + user

	existsInBloomFilterCmd := redis.NewIntCmd(redisBloomFilterExistsCmd, redisBloomFilterName, article)
	recoUtil.redisClient.Process(existsInBloomFilterCmd)
	cmdResult, _ := existsInBloomFilterCmd.Result()

	return int(cmdResult)
}

//CloseConn ...
func (recoUtil *RecommendationUtil) CloseConn() {
	recoUtil.redisClient.Close()
}
