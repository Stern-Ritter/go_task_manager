package app

import (
	"flag"

	"github.com/caarlos0/env"

	"github.com/Stern-Ritter/go_task_manager/internal/config"
)

func GetConfig(c config.ServerConfig) (config.ServerConfig, error) {
	parseFlags(&c)

	err := env.Parse(&c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func parseFlags(c *config.ServerConfig) {
	flag.IntVar(&c.Port, "p", 7540, "port to run server")
	flag.StringVar(&c.DatabaseFile, "f", "scheduler.db", "database file name")
	flag.Parse()
}
