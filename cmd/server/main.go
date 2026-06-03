package main

import (
	logger "core/pkg/log"
	router "core/api/router"
	scheduler "core/pkg/schedule"
	config "core/internal/config"
	database "core/internal/database"

	"net/http"
	"time"

	"github.com/ory/graceful"
)

// @title SendIn API
// @version 1.0
// @description LinkedIn outreach automation API
// @BasePath /api/v1

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name sendin_auth
func main() {
	logger.Start()

	logger.Info("🧩 Configuring environment")
	cfg, err := config.Load()
	logger.Fatal(err, "Failed to load env")

	logger.Info("📡 Bootstrapping database")
	gormDB := database.Init(&cfg.Database)
	sqlDB, err := gormDB.DB()
	logger.Fatal(err, "Failed to create connection with database")
	defer sqlDB.Close()

	server := graceful.WithDefaults(
		&http.Server{
			Addr:         cfg.Server.Port,
			Handler:      router.Config(gormDB, &cfg.Server),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	)

	err = database.Migrate(&cfg.Database)
	logger.Fatal(err, "Failed to create migrations")

	logger.Info("⏰ Scheduler started (polling every 1 min)")
	go scheduler.Start(gormDB)

	logger.Info("🚀 Starting server on port %s", cfg.Server.Port)
	err = graceful.Graceful(
		server.ListenAndServe, 
		server.Shutdown,
	)

	logger.Fatal(err, "Failed to gracefully shutdown")
	logger.Info("🛑 Server exited")
}

