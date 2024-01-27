package routes

import (
	"github.com/fvdime/keen-go-backend/controllers"
	"github.com/fvdime/keen-go-backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(routes *gin.Engine) {
	routes.Use(middleware.Authenticate())
	routes.GET("/user/:user_id", controllers.GetUser())
	// routes.GET("/users/:user_id", controllers.GetUser())
}
