package main

import (
	"github.com/fvdime/keen-go-backend/routes"

	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"SUCCESS": "Access granted!"})
	})

	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"SUCCESS": "Access granted for api-2!"})
	})

	router.Run(":" + port)
}
