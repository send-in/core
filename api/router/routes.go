package router

import (
	controller "core/api/controller"
	middleware "core/api/middleware"
	config "core/internal/config"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "core/docs"
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
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
		},
	))

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
		),
	)

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
		protected.DELETE("/messages/:id", controllers.DeleteMessage)

		// Templates
		protected.GET("/templates", controllers.GetTemplates)
		protected.GET("/templates/:id", controllers.GetTemplate)
		protected.POST("/templates", controllers.CreateTemplate)
		protected.PUT("/templates/:id", controllers.UpdateTemplate)
		protected.DELETE("/templates/:id", controllers.DeleteTemplate)
	}

	return router
}