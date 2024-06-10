package main

import (
	"fmt"
	_ "net/http/pprof"

	"github.com/rutkin/url-shortener/internal/app"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/logger"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n Build date: %s\n Build commit: %s\n", buildVersion, buildDate, buildCommit)

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
	defer server.Close()
}
