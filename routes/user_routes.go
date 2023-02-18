package routes

import (
	"github.com/demola234/real-estate/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
	incomingRoutes.POST("/users/:user_id", controllers.LogOut())
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/users/verifyOtp", controllers.VerifyOtp())
	incomingRoutes.POST("/users/resendOtp", controllers.ResendOtp())
	incomingRoutes.POST("/users/createPassword", controllers.CreatePassword())
	incomingRoutes.POST("/users/forgetPassword", controllers.ForgotPassword())
	incomingRoutes.POST("/users/resetPassword", controllers.ResetPassword())
	incomingRoutes.POST("/users/google", controllers.GoogleOauth())
	incomingRoutes.POST("/users/facebook", controllers.FacebookOauth())
	incomingRoutes.PATCH("/users/updateImage", controllers.UpdateProfileImage())
	incomingRoutes.DELETE("/users/:user_id", controllers.DeleteUser())
}
