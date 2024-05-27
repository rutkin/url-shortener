// package config, contains settings for server
package config

import (
	"flag"
	"fmt"
	"os"
)

// NetAddress - type for network address
type NetAddress string

// Config - configuration type
type Config struct {
	Server          NetAddress
	Base            NetAddress
	LogLevel        string
	FileStoragePath string
	DatabaseDSN     string
}

// ServerConfig - default server settings, address - http://localhost:8080, log level - info, storage - file
var ServerConfig = Config{Server: "localhost:8080", Base: "http://localhost:8080", LogLevel: "info", FileStoragePath: "/tmp/short-url-db.json"}

// return network address string
func (a NetAddress) String() string {
	return string(a)
}

// set network address
func (a *NetAddress) Set(s string) error {
	*a = NetAddress(s)
	return nil
}

// parse config from argument and environment variables
func ParseFlags() error {
	flag.Var(&ServerConfig.Server, "a", "http server address")
	flag.Var(&ServerConfig.Base, "b", "base server address")
	flag.StringVar(&ServerConfig.LogLevel, "l", "info", "log level")
	flag.StringVar(&ServerConfig.FileStoragePath, "f", "/tmp/short-url-db.json", "file storage path")
	flag.StringVar(&ServerConfig.DatabaseDSN, "d", "", "database dsn")
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

	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		ServerConfig.FileStoragePath = fileStoragePath
	}

	if databaseDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		ServerConfig.DatabaseDSN = databaseDSN
	}

	return nil
}
