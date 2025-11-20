package main

import (
	"log"
	"os"

	"contactsAI/contacts/docs"
	"contactsAI/contacts/internal/config"
	"contactsAI/contacts/internal/routing"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Contacts AI API
//	@version		1.0
//	@description	API for managing contacts with AI features
//	@host			localhost:33500
//	@contact.name	Mateusz Kurowski
//	@contact.email	mateusz.kurowski28@gmail.com
//	@BasePath		/api.

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		if loadErr := godotenv.Load("./env/.env"); loadErr != nil {
			panic("Error loading .env file")
		}
	}

	env, connErr := config.NewEnv(os.Getenv("DB_URL"), false)
	if env.Logger == nil {
		log.Println("Failed to setup logger.")
	}
	if connErr != nil || env.Queries == nil {
		env.Logger.Error("Failed to initialize environment", "error", connErr)
	}
	if env.Bucket == nil {
		env.Logger.Error("Failed to initialize s3 connection", "error", connErr)
	}

	router := routing.SetupRouter(env)

	// Swagger
	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// PORT is set via environment variable
	// Default to 8080 if not set
	if runErr := router.Run(); runErr != nil {
		env.Logger.Error("Failed to start server", "error", runErr)
	}
}
