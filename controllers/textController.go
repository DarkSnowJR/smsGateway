package controllers

import (
	"smsGateway/initializers"
	"smsGateway/models"

	"github.com/gin-gonic/gin"
)

func TextList(c *gin.Context) {
	// Get the message ID from the request parameters
	messageID := c.Param("messageID")

	// Get the message from the database
	var message models.Message
	if err := initializers.DB.First(&message, messageID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Message not found"})
		return
	}

	var texts []models.Text
	if err := initializers.DB.Where("message_id = ?", messageID).Find(&texts).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"texts": texts,
	})
}
