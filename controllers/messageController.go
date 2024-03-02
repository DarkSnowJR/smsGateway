package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"smsGateway/initializers"
	"smsGateway/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Define a threshold for very high and very low sending rates
const highSendingRateThreshold = 10000 // You can adjust this threshold based on your requirements
const lowSendingRateThreshold = 1000    // You can adjust this threshold based on your requirements

func MessageCreate(c *gin.Context) {
	var body struct {
		To      []string `json:"to" binding:"required"`
		Content string   `json:"content" binding:"required"`
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

	// Marshal "To" field to JSON
	toJSON, err := json.Marshal(body.To)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	message := models.Message{
		To:      []string{string(toJSON)}, // Wrap in a slice to ensure it's a JSON array
		Content: body.Content,
		UserID:  user.ID,
	}

	if err := initializers.DB.Create(&message).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Send text message
	go sendTextTask(message.ID)

	c.JSON(200, gin.H{
		"message": message,
	})
}

func MessageList(c *gin.Context) {
	var user models.User
	if err := initializers.DB.Where("user_name = ?", c.Param("username")).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var messages []models.Message
	if err := initializers.DB.Where("user_id = ?", user.ID).Find(&messages).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"messages": messages,
	})
}

func sendTextTask(messageID uint) {
	fmt.Println("Text task Started for message ID:", messageID)

	// Declare a variable to store the message
	var message models.Message

	// Fetch the message from the database
	if err := initializers.DB.First(&message, messageID).Error; err != nil {
		fmt.Println("Error fetching message:", err)
		return
	}

	// Check if the To field is populated in the retrieved message
	if len(message.To) == 0 {
		fmt.Println("Message.To is empty. Cannot proceed.")
		return
	}

	// Fetch the user associated with the message
	var user models.User
	if err := initializers.DB.First(&user, message.UserID).Error; err != nil {
		fmt.Println("Error fetching user:", err)
		return
	}

	// Fetch the number of text messages for the user
	var textMessageCount int64
	if err := initializers.DB.Model(&models.Text{}).
		Joins("JOIN messages ON texts.message_id = messages.id").
		Where("messages.user_id = ?", user.ID).
		Count(&textMessageCount).Error; err != nil {
		fmt.Println("Error fetching text message count for user:", err)
		return
	}

	// Calculate the sending rate based on the number of text messages
	sendingRate := float64(textMessageCount)

	fmt.Println("Sending Rate: ", sendingRate)

	// Check the sending rate and apply different logic based on it
	if sendingRate >= highSendingRateThreshold {
		// Logic for very high sending rate customers
		fmt.Println("Handling very high sending rate...")
	} else if sendingRate <= lowSendingRateThreshold {
		// Logic for very low sending rate customers
		fmt.Println("Handling very low sending rate...")
	} else {
		// Default logic for normal sending rate customers
		fmt.Println("Handling normal sending rate...")
	}

	// Get the message price from the environment variable
	messagePriceStr := os.Getenv("MESSAGE_PRICE")
	messagePrice, err := strconv.ParseFloat(messagePriceStr, 64)
	if err != nil {
		fmt.Println("Error parsing MESSAGE_PRICE:", err)
		return
	}

	// Decode the JSON string to get the array of phone numbers
	var receivers []string
	if err := json.Unmarshal([]byte(message.To[0]), &receivers); err != nil {
		fmt.Println("Error decoding receivers:", err)
		return
	}

	// Check if the user's balance is below the limit
	if user.Balance < messagePrice {
		fmt.Println("User's balance is below the limit. Cannot send texts.")
		return
	}

	// Loop through receivers and create Text records
	for _, receiver := range receivers {

		// Check if the user's balance is sufficient for this text
		if user.Balance >= messagePrice {
			text := models.Text{
				To:        receiver,
				MessageID: messageID,
				Content:   message.Content,
				Status:    true,
			}
			initializers.DB.Create(&text)

			tx := initializers.DB.Begin()
			// Deduct cost and update user balance
			user.Balance -= messagePrice
			if err := tx.Save(&user).Error; err != nil {
				tx.Rollback()
				// Handle the error, e.g., log it or return an error response
				fmt.Println("Error updating user balance:", err)
				return
			}

			// Commit the transaction
			tx.Commit()

		} else {
			// If balance is not enough, set text status to false
			text := models.Text{
				To:        receiver,
				MessageID: messageID,
				Content:   message.Content,
				Status:    false,
			}
			initializers.DB.Create(&text)
		}
	}

	// Update the status of the message
	message.Status = true
	if err := initializers.DB.Save(&message).Error; err != nil {
		fmt.Println("Error updating message status:", err)
		return
	}

	// Update the user's balance
	if err := initializers.DB.Save(&user).Error; err != nil {
		fmt.Println("Error updating user balance:", err)
		return
	}

	fmt.Println("Text task Completed for message ID:", messageID)
}
