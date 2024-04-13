package main

import (
	"log"

	"github.com/Stern-Ritter/go_task_manager/internal/app"
	"github.com/Stern-Ritter/go_task_manager/internal/config"
	"github.com/Stern-Ritter/go_task_manager/internal/logger"
	"go.uber.org/zap"
)

func main() {
	config, err := app.GetConfig(config.ServerConfig{
		DatabaseDriverName: "sqlite",
		LoggerLvl:          "info",
	})
	if err != nil {
		log.Fatalf("%+v", err)
		return
	}

	logger, err := logger.Initialize(config.LoggerLvl)
	if err != nil {
		log.Fatalf("%+v", err)
		return
	}

	err = app.Run(&config, logger)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "start server"))
	}
}
