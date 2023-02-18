package main

import (
	"os"

	"github.com/demola234/real-estate/middleware"
	"github.com/demola234/real-estate/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	routes.PropertiesRoutes(router)
	routes.PaymentRoutes(router)
	routes.NotificationRoutes(router)
	routes.ManageUserRoutes(router)
	routes.AgentRoutes(router)

	router.Run(":" + port)
}
