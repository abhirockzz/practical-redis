package main

import (
	"api"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(api.RedisConnHandler)

	router.POST("users", api.UserRegistrationHandler)

	router.POST("news", api.PostNewsItemHandler)
	router.GET("news/:newsid", api.GetNewsItemByIDHandler)
	router.GET("news", api.GetAllNewsItemsHandler)

	router.GET("news/:newsid/comments", api.GetCommentsForNewsID)
	router.POST("news/:newsid/comments", api.PostCommentsForNewsItemHandler)

	router.GET("news/:newsid/upvotes", api.GetUpvotesForNewsItemHandler)
	router.POST("news/:newsid/upvotes", api.UpvoteNewsItemHandler)

	router.Run()
}
