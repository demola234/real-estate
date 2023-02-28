package routes

import (
	"github.com/demola234/real-estate/controllers"
	"github.com/gin-gonic/gin"
)

func NotificationRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/notifications", controllers.GetNotifications())
	incomingRoutes.GET("/notifications/:notification_id", controllers.GetNotification())
	incomingRoutes.DELETE("/notifications/:notification_id", controllers.DeleteNotification())
	incomingRoutes.POST("/notifications", controllers.TestPushNotification())
}
