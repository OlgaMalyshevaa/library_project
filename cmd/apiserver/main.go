package main

import (
	_ "library_project/cmd/apiserver/docs"
	"library_project/database"
	"library_project/internal/middleware"
	"library_project/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	middleware.SetupLogger()

	database.ConnectDatabase()

	r := gin.Default()

	routes.SetupRoutes(r)

	r.GET("/swagger/*any", func(c *gin.Context) {
		log.Println("Swagger UI requested")
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
