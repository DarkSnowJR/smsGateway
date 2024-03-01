package controllers

import (
	"smsGateway/initializers"
	"smsGateway/models"

	"github.com/gin-gonic/gin"
)

func UserCreate(c *gin.Context) {
	// Get data of request body
	var body struct {
		UserName string
		Name     string
		Email    string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user := models.User{
		UserName: body.UserName,
		Name:     body.Name,
		Email:    body.Email,
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}

func GetBalance(c *gin.Context) {
	var user models.User
	if err := initializers.DB.Where("user_name = ?", c.Param("username")).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, gin.H{
		"balance": user.Balance,
	})
}

func AddBalance(c *gin.Context) {
	var body struct {
		Balance float64
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := initializers.DB.Where("user_name = ?", c.Param("username")).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := initializers.DB.Model(&user).Update("balance", user.Balance+body.Balance).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"balance": user.Balance,
	})
}
