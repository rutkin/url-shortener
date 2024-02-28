package main

import (
	"github.com/rutkin/url-shortener/internal/app"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/logger"
)

func main() {
	err := config.ParseFlags()

	if err != nil {
		panic(err)
	}

	err = logger.Initialize(config.ServerConfig.LogLevel)

	if err != nil {
		panic(err)
	}

	server, err := app.NewServer()
	if err != nil {
		panic(err)
	}

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
