package util

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/go-redis/redis"
)

//LoadArticlesForTopics - add articles for each topic
//e.g. SET topic:software:articles , topic:programming:articles
func LoadArticlesForTopics(redisCoordinate string) {

	client := redis.NewClient(&redis.Options{Addr: redisCoordinate})
	_, err := client.Ping().Result()

	if err != nil {
		fmt.Println("Could not connect to Redis " + err.Error())
	} else {
		fmt.Println("Connected to Redis " + redisCoordinate)
	}
	defer client.Close()

	//software engineering
	client.SAdd("topic:softwareengineering:articles", "https://medium.com/@anildash/what-if-javascript-wins-84898e5341a")
	client.SAdd("topic:softwareengineering:articles", "https://hackernoon.com/the-7-biggest-lessons-ive-learned-by-building-a-twitter-bot-59fee84a9ed9")
	client.SAdd("topic:softwareengineering:articles", "https://towardsdatascience.com/data-science-for-startups-data-pipelines-786f6746a59a")
	client.SAdd("topic:softwareengineering:articles", "https://towardsdatascience.com/universal-language-model-to-boost-your-nlp-models-d59469dcbd64")
	client.SAdd("topic:softwareengineering:articles", "https://towardsdatascience.com/designing-an-iot-solution-in-2018-7fe1356e63d6")

	fmt.Println("loaded softengg topics")

	//creativity
	client.SAdd("topic:creativity:articles", "https://medium.com/swlh/how-to-make-something-people-love-a8364771b7e6")
	client.SAdd("topic:creativity:articles", "https://medium.com/personal-growth/walt-disney-how-to-truly-love-what-you-do-f3449c78ca65")
	client.SAdd("topic:creativity:articles", "https://medium.com/@michaelpollan/medium-com-trips-aed86f968810")
	client.SAdd("topic:creativity:articles", "https://blog.prototypr.io/growing-an-idea-from-an-interest-to-a-product-a0757b415bbb")
	client.SAdd("topic:creativity:articles", "https://medium.com/@shauntagrimes/challenge-yourself-to-learn-from-masters-3f99064e0f2e")

	fmt.Println("loaded creativity topics")

	//programming
	client.SAdd("topic:programming:articles", "https://towardsdatascience.com/unsupervised-learning-with-python-173c51dc7f03")
	client.SAdd("topic:programming:articles", "https://medium.com/@evheniybystrov/react-redux-for-lazy-developers-b551f16a456f")
	client.SAdd("topic:programming:articles", "https://medium.com/hackerpreneur-magazine/how-i-hacked-into-one-of-the-most-popular-dating-websites-4cb7907c3796")
	client.SAdd("topic:programming:articles", "https://medium.com/sololearn/warning-your-programming-career-b9579b3a878b")
	client.SAdd("topic:programming:articles", "https://medium.com/@jrodthoughts/using-deep-learning-to-understand-your-source-code-28e5c284bfda")

	fmt.Println("loaded programming topics")

	//productivity
	client.SAdd("topic:productivity:articles", "https://medium.com/thrive-global/one-battle-you-need-to-win-everyday-a3635a5562f")
	client.SAdd("topic:productivity:articles", "https://medium.com/swlh/try-this-sprint-approach-to-work-if-you-want-to-do-more-in-less-time-fd15fd634b22")
	client.SAdd("topic:productivity:articles", "https://medium.com/swlh/10-daily-habits-of-the-most-successful-entrepreneurs-9a0bb5e9e91b")
	client.SAdd("topic:productivity:articles", "https://theascent.pub/this-is-how-playing-musical-chairs-can-make-you-more-productive-6793893d9e84")
	client.SAdd("topic:productivity:articles", "https://medium.com/the-polymath-project/you-are-not-your-goals-446559f1f118")

	fmt.Println("loaded productivity topics")

	//travel
	client.SAdd("topic:travel:articles", "https://medium.com/@krisgage/why-im-over-airbnb-f4a35aacc951")
	client.SAdd("topic:travel:articles", "https://medium.com/@jeffgoins/3-ways-your-life-will-change-if-you-travel-while-youre-young-c66d5a1993ee")
	client.SAdd("topic:travel:articles", "https://medium.com/@sravss/of-the-merlion-and-the-bao-and-tom-yum-fish-soup-2eec2ff5cf82")
	client.SAdd("topic:travel:articles", "https://medium.com/andrew-across-america/fake-conversations-cb87e50cd58")
	client.SAdd("topic:travel:articles", "https://medium.com/@sfarah214/what-do-i-do-with-the-memories-b3c12e446d94")

	fmt.Println("loaded travel topics")

	fmt.Println("closed redis connection...")
}

//LoadUserInterests - add topics whic users are interested in
//e.g. SET user:abhi:interests = {software, programming}
func LoadUserInterests(redisCoordinate string) {
	topics := []string{"softwareengineering", "creativity", "programming", "productivity", "travel"}

	client := redis.NewClient(&redis.Options{Addr: redisCoordinate})
	_, err := client.Ping().Result()

	if err != nil {
		fmt.Println("Could not connect to Redis " + err.Error())
	} else {
		fmt.Println("Connected to Redis " + redisCoordinate)
	}
	defer client.Close()

	rand.Seed(50)

	for i := 1; i <= 4; i++ {
		setName := "user:user-" + strconv.Itoa(i) + ":interests"

		//try to add (max) 5 interests per user. not all might be added because
		//we are at the mercy of the random generator
		for c := 0; c < 5; c++ {
			topic := topics[rand.Intn(len(topics))]
			result, _ := client.SAdd(setName, topic).Result()
			if result > 0 {
				fmt.Println("added topic " + topic + " to set " + setName)
			}
		}

	}

	fmt.Println("closed redis connection...")

}
