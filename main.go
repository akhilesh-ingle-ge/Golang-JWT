package main

import (
	"log"
	"os"

	"github.com/akhilesh-ingle-ge/jwt/config"
	"github.com/akhilesh-ingle-ge/jwt/controllers"
	"github.com/akhilesh-ingle-ge/jwt/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error in loading .env file")
	}
	config.SetUpDB()
}

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Message": "Hello World",
		})
	})

	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)

	router.Run(":" + os.Getenv("PORT"))
}
