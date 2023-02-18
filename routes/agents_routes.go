package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/demola234/real-estate/controllers"
)

func AgentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/agents", controllers.GetAgents())
	incomingRoutes.GET("/agents/:agent_id", controllers.GetAgent())
	incomingRoutes.POST("/agents", controllers.BecomeAnAgent())
	
	incomingRoutes.PUT("/agents/:agent_id", controllers.UpdateAgent())
	incomingRoutes.DELETE("/agents/:agent_id", controllers.DeleteAgent())
}
