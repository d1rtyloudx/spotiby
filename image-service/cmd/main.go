package main

import (
	"github.com/d1rtyloudx/spotiby-pkg/logger"
	"github.com/joho/godotenv"
	"image-service/internal/app"
	"image-service/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("cannot load .env file: " + err.Error())
	}

	cfg := config.MustLoad()

	log := logger.New()

	application := app.NewApp(log, cfg)

	err := application.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
