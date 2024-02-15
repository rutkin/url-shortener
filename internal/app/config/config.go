package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type NetAddress string

type Config struct {
	Server NetAddress
	Base   NetAddress
}

var ServerConfig = Config{Server: "localhost:8080", Base: "http://localhost:8080"}
var errInvalidAddress = errors.New("need address in a form host:port")

func (a NetAddress) String() string {
	return string(a)
}

func (address *NetAddress) Set(s string) error {
	*address = NetAddress(s)
	return nil
}

func init() {
	flag.Var(&ServerConfig.Server, "a", "http server address")
	flag.Var(&ServerConfig.Base, "b", "base server address")
}

func (c Config) ParseFlags() error {
	flag.Parse()

	if serverAddress, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		err := c.Server.Set(serverAddress)
		if err != nil {
			return fmt.Errorf("Failed to set server address '%s' in config", serverAddress)
		}
	}

	if baseAddress, ok := os.LookupEnv("BASE_ADDRESS"); ok {
		err := c.Server.Set(baseAddress)
		if err != nil {
			return fmt.Errorf("Failed to set base address '%s' in config", baseAddress)
		}
	}

	return nil
}
