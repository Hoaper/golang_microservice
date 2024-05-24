package main

import (
	"EffectiveMobile/config"
	"EffectiveMobile/internal/app"
	"go.uber.org/zap"
)

func main() {

	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	var logger *zap.Logger

	if cfg.DEBUG {
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	}

	newApp, err := app.NewApp(logger, cfg)
	if err != nil {
		logger.Fatal("Failed to create newApp")
	}
	newApp.Run()

}
