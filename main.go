package main

import (
	"log"
	"os"

	"contactsAI/contacts/docs"
	"contactsAI/contacts/internal/config"
	"contactsAI/contacts/internal/routing"

	"github.com/gin-contrib/sessions"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			Contacts AI API
// @version		1.0
// @description	API for managing contacts with AI features
// @host			localhost:33500
// @contact.name	Mateusz Kurowski
// @contact.email	mateusz.kurowski28@gmail.com
// @BasePath		/api.
func main() {
	if os.Getenv("GIN_MODE") != "release" {
		if loadErr := godotenv.Load("./env/.env"); loadErr != nil {
			panic("Error loading .env file")
		}
	}

	env, connErr := config.NewEnv(os.Getenv("DB_URL"))
	if connErr != nil || env.Queries == nil {
		log.Fatalf("Failed to initialize database: %v", connErr)
	}

	// Setup authentication providers
	config.SetupProviders()

	router := routing.SetupRouter(env)

	router.Use(sessions.Sessions("mysession", env.CookieStore))

	routing.RegisterRoutes(router, env)

	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// PORT is set via environment variable
	// Default to 8080 if not set
	if runErr := router.Run(); runErr != nil {
		log.Fatalf("Failed to start server: %v", runErr)
	}
}
