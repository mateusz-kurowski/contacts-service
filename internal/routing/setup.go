package routing

import (
	"contactsAI/contacts/internal/config"
	"contactsAI/contacts/internal/handlers"
	"contactsAI/contacts/internal/middleware"
	"contactsAI/contacts/internal/validation"

	"github.com/gin-gonic/gin"
)

func SetupRouter(env *config.Env) *gin.Engine {
	router := gin.Default()

	middleware.SetupMiddlewares(router)

	validation.SetupValidation()
	// Register routes so callers that only call SetupRouter
	// (for example tests) get a router with all endpoints wired.
	RegisterRoutes(router, env)

	return router
}

func RegisterRoutes(router *gin.Engine, env *config.Env) {
	apiGroup := router.Group("/api")

	// Register routes
	handlers.RegisterContactsRoutes(apiGroup, env)
}
