package lcm

import (
	"fmt"
	"service"

	"github.com/gin-gonic/gin"
)

func StartServiceHandler(c *gin.Context) {
	fmt.Println("StartServiceHandler API invoked")
	if service.GetTweetsListenerStatus() {
		alreadyRunning := "Tweets listener service is already running!"
		fmt.Println(alreadyRunning)
		c.Writer.WriteString(alreadyRunning)
		return
	}
	err := service.Start()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Writer.WriteString("Started Tweets listener")
}

func StopServiceHandler(c *gin.Context) {
	fmt.Println("StopServiceHandler API invoked")
	if !service.GetTweetsListenerStatus() {
		notRunning := "Tweets listener service is not running!"
		fmt.Println(notRunning)
		c.Writer.WriteString(notRunning)
		return
	}
	service.Stop()
	c.Writer.WriteString("Stopped Tweets listener")
}
