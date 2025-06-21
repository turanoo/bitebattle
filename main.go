package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/turanoo/bitebattle/api"
	"github.com/turanoo/bitebattle/internal/auth"
	"github.com/turanoo/bitebattle/pkg/config"
	"github.com/turanoo/bitebattle/pkg/db"
	"github.com/turanoo/bitebattle/pkg/logger"
	"github.com/turanoo/bitebattle/pkg/utils"
)

var (
	MIGRATIONS_PATH = "./migrations"
)

func main() {
	logger.Init()

	ctx := context.Background()
	cfg, err := config.LoadConfig(ctx, "config")
	if err != nil {
		logger.Errorf("Failed to load config: %v", err)
		os.Exit(1)
	}
	auth.InitJWTKey(cfg)

	if err := db.Init(cfg); err != nil {
		logger.Errorf("Failed to connect to DB: %v", err)
		os.Exit(1)
	}
	database := db.GetDB()

	migrationLogger := logger.Log.WithFields(logrus.Fields{"requestId": "startup"})
	if err := utils.RunMigrations(db.GetPostgresURL(cfg), MIGRATIONS_PATH, migrationLogger); err != nil {
		logger.Errorf("Failed to run migrations: %v", err)
		os.Exit(1)
	}

	gin.SetMode(cfg.Gin.Mode)
	switch cfg.Gin.Log.Level {
	case "debug":
		logger.Log.SetLevel(logrus.DebugLevel)
	case "info":
		logger.Log.SetLevel(logrus.InfoLevel)
	default:
		logger.Log.SetLevel(logrus.InfoLevel)
	}

	if cfg.Gin.Log.Format == "json" {
		logger.Log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.Log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	router := gin.New()
	router.Use(logger.Middleware())
	router.Use(logger.RequestLogger())
	router.Use(logger.ErrorRecovery())

	api.SetupRoutes(router, database, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Infof("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Errorf("Server failed: %v", err)
		os.Exit(1)
	}
}
