package config

import (
	"flag"
	"fmt"
	"os"
)

type NetAddress string

type Config struct {
	Server   NetAddress
	Base     NetAddress
	LogLevel string
}

var ServerConfig = Config{Server: "localhost:8080", Base: "http://localhost:8080", LogLevel: "info"}

func (a NetAddress) String() string {
	return string(a)
}

func (a *NetAddress) Set(s string) error {
	*a = NetAddress(s)
	return nil
}

func ParseFlags() error {
	flag.Var(&ServerConfig.Server, "a", "http server address")
	flag.Var(&ServerConfig.Base, "b", "base server address")
	flag.StringVar(&ServerConfig.LogLevel, "l", "info", "log level")
	flag.Parse()

	if serverAddress, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		err := ServerConfig.Server.Set(serverAddress)
		if err != nil {
			return fmt.Errorf("failed to set server address '%s' in config", serverAddress)
		}
	}

	if baseAddress, ok := os.LookupEnv("BASE_ADDRESS"); ok {
		err := ServerConfig.Server.Set(baseAddress)
		if err != nil {
			return fmt.Errorf("failed to set base address '%s' in config", baseAddress)
		}
	}

	return nil
}
