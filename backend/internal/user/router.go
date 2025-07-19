package user

import (
	"document-reader-chatbot/configs"
	"document-reader-chatbot/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all user-related routes
func SetupRoutes(router *gin.RouterGroup, service Service) {
	controller := NewController(service)

	// Load configuration for JWT middleware
	cfg, err := configs.Load()
	if err != nil {
		panic("Failed to load configuration for user routes")
	}

	// Public routes (no authentication required)
	auth := router.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/register", controller.CreateUser)
	}

	// Protected routes for user profile management
	profile := router.Group("/profile")
	profile.Use(middleware.JWTAuth(cfg.JWT))
	{
		profile.GET("", controller.GetProfile)
		profile.PUT("", controller.UpdateProfile)
	}

	// Admin routes for user management
	users := router.Group("/users")
	users.Use(middleware.JWTAuth(cfg.JWT))
	{
		// List users
		users.GET("", controller.ListUsers)

		// Get specific user
		users.GET("/:id", controller.GetUser)

		// Update user
		users.PUT("/:id", controller.UpdateUser)

		// Delete user
		users.DELETE("/:id", controller.DeleteUser)

		// Change password
		users.POST("/:id/change-password", controller.ChangePassword)
	}
}
