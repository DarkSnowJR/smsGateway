package main

import (
	"smsGateway/initializers"
	"smsGateway/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Message{})
	initializers.DB.AutoMigrate(&models.Text{})
}
