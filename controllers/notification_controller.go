package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/demola234/real-estate/interfaces"
	"github.com/demola234/real-estate/services"
	"github.com/gin-gonic/gin"
)

func GetNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}

func GetNotifications() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}

func CreateNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}

func DeleteNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}

func SendushNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}

func SendEmailNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}

func TestPushNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var push interfaces.PushNotificationToUser

		if err := c.BindJSON(&push); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := services.SendPushNotificationToUser(&push)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

}
