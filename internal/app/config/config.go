// package config, contains settings for server
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// NetAddress - type for network address
type NetAddress string

// Config - configuration type
type Config struct {
	Server          NetAddress `json:"server_address"`
	Base            NetAddress `json:"base_url"`
	LogLevel        string
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
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
	var configPath string
	flagServerConfig := ServerConfig
	flag.StringVar(&configPath, "c", "", "config file path")
	flag.StringVar(&configPath, "config", "", "config file path")
	flag.Var(&flagServerConfig.Server, "a", "http server address")
	flag.Var(&flagServerConfig.Base, "b", "base server address")
	flag.StringVar(&flagServerConfig.LogLevel, "l", "info", "log level")
	flag.StringVar(&flagServerConfig.FileStoragePath, "f", "/tmp/short-url-db.json", "file storage path")
	flag.StringVar(&flagServerConfig.DatabaseDSN, "d", "", "database dsn")
	flag.BoolVar(&flagServerConfig.EnableHTTPS, "s", false, "enable https")
	flag.StringVar(&flagServerConfig.TrustedSubnet, "t", "", "trusted subnet")
	flag.Parse()

	if len(configPath) > 0 {
		config, err := os.Open(configPath)
		if err != nil {
			return err
		}
		if err := json.NewDecoder(config).Decode(&ServerConfig); err != nil {
			return err
		}
	}

	ServerConfig = flagServerConfig

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

	if enableHTTPS, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		var err error
		ServerConfig.EnableHTTPS, err = strconv.ParseBool(enableHTTPS)
		if err != nil {
			return fmt.Errorf("failed to parse enable https bool value from '%s'", enableHTTPS)
		}
	}

	if trustedSubnet, ok := os.LookupEnv("TRUSTED_SUBNET"); ok {
		ServerConfig.TrustedSubnet = trustedSubnet
	}

	return nil
}
