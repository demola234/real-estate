package routes

import (
	"github.com/demola234/real-estate/controllers"
	"github.com/gin-gonic/gin"
)

func ManageUserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/profile", controllers.GetProfile())
	incomingRoutes.PUT("/profile", controllers.UpdateProfile())
	incomingRoutes.PUT("/profile/password", controllers.ChangePassword())
}
