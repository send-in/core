package router

import (
	controller "core/api/controller"
	middleware "core/api/middleware"
	config "core/internal/config"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Config(
	db *gorm.DB,
	cfg *config.ServerConfig,
) http.Handler {
	router := gin.Default()
	controllers := controller.Create(db)

	router.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				"GET",
				"POST",
				"PUT",
				"DELETE",
			},
			AllowHeaders: []string{"*"},
			AllowCredentials: true,
		},
	))

	v1 := router.Group("/api/v1")

	// Health
	v1.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"status": "healthy",
				},
			)
		},
	)

	// Public Routes
	auth := v1.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/signup", controllers.Signup)
	}

	// Protected Routes
	protected := v1.Group("/")
	protected.Use(middleware.Authenticate(db))
	{
		// Account
		protected.GET("/account", controllers.GetAccount)
		protected.PUT("/account", controllers.UpdateAccount)
		protected.DELETE("/account", controllers.DeleteAccount)

		// Connections
		protected.GET("/connections", controllers.GetConnections)
		protected.POST("/connections", controllers.CreateConnections)

		// Messages
		protected.GET("/messages", controllers.GetMessages)
		protected.GET("/messages/:id", controllers.GetMessage)
		protected.POST("/messages", controllers.CreateMessage)
		protected.DELETE("/messages/:id", controllers.DeleteMessage,)
	}

	return router
}