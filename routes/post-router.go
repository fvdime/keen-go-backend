package routes

import (
	"github.com/fvdime/keen-go-backend/controllers"

	"github.com/gin-gonic/gin"
)

func PostRoutes(routes *gin.Engine) {
	routes.POST("post/create", controllers.CreatePost())
	routes.DELETE("post/:post_id", controllers.DeletePost())
	routes.PATCH("post/:post_id", controllers.UpdatePost())
	routes.GET("post/:post_id", controllers.GetPost())
	routes.GET("post/", controllers.GetPosts())
}
