package routes

import (
	"github.com/fvdime/keen-go-backend/controllers"

	"github.com/gin-gonic/gin"
)


func AuthRoutes(routes *gin.Engine){
	routes.POST("auth/signup", controllers.SignUp())
	routes.POST("auth/signin", controllers.SignIn())
}