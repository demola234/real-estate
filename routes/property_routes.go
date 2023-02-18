package routes

import (
	"github.com/demola234/real-estate/controllers"
	"github.com/gin-gonic/gin"
)

func PropertiesRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/properties", controllers.GetProperties())
	incomingRoutes.GET("/properties/:property_id", controllers.GetProperty())
	incomingRoutes.POST("/properties", controllers.AddProperty())
	incomingRoutes.PATCH("/properties/:property_id", controllers.UpdateProperty())
	incomingRoutes.DELETE("/properties/:property_id", controllers.DeleteProperty())
	incomingRoutes.GET("/properties/agent/:agent_id", controllers.GetPropertyByAgent())
	incomingRoutes.GET("/properties/location/:location", controllers.GetPropertyByLocation())
	incomingRoutes.GET("/properties/featured", controllers.FeaturedProperties())
	incomingRoutes.POST("/properties/rate", controllers.RateProperty())
}
