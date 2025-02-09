package main

import (
	"github.com/d1rtyloudx/spotiby-pkg/logger"
	"github.com/joho/godotenv"
	"user-service/internal/app"
	"user-service/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	cfg := config.MustLoad()

	log := logger.New()

	application := app.NewApp(log, cfg)
	if err := application.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
