package main

import (
	router "core/api/router"
	config "core/internal/config"
	"core/internal/database"
	logger "core/pkg/log"

	"net/http"
	"time"

	"github.com/ory/graceful"
)

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

	err = database.Migrate(gormDB)
	logger.Fatal(err, "Failed to create migrations")

	logger.Info("🚀 Starting server on port %s", cfg.Server.Port)
	err = graceful.Graceful(
		server.ListenAndServe, 
		server.Shutdown,
	)

	logger.Fatal(err, "Failed to gracefully shutdown")
	logger.Info("🛑 Server exited")
}

