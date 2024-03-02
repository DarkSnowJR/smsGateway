package main

import (
	"smsGateway/controllers"
	"smsGateway/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()

	// Routes
	router.POST("/users", controllers.UserCreate)
	router.GET("/users/:username/balance", controllers.GetBalance)
	router.POST("/users/:username/balance", controllers.AddBalance)

	router.POST("/messages/:username", controllers.MessageCreate)
	router.GET("/messages/:username", controllers.MessageList)

	router.GET("/texts/:messageID", controllers.TextList)

	err := router.Run()
	if err != nil {
		return
	}
}
