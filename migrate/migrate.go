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
	err := initializers.DB.AutoMigrate(&models.User{})
	if err != nil {
		return
	}
	err = initializers.DB.AutoMigrate(&models.Message{})
	if err != nil {
		return
	}
	err = initializers.DB.AutoMigrate(&models.Text{})
	if err != nil {
		return
	}
}
