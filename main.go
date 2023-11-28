package main

import (
	"github.com/gin-gonic/gin"
	"github.com/marsh-sudo/go-jwt/controllers"
	"github.com/marsh-sudo/go-jwt/initializers"
	"github.com/marsh-sudo/go-jwt/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.SignUp) 
	r.POST("/login",controllers.Login)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.Run() // listen and serve on 0.0.0.0:8080
}
