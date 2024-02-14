package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type NetAddress struct {
	Scheme string
	Host   string
	Port   int
}

type Config struct {
	Server NetAddress
	Base   NetAddress
}

var ServerConfig = Config{Server: NetAddress{"", "localhost", 8080}, Base: NetAddress{"http", "localhost", 8080}}
var errInvalidAddress = errors.New("need address in a form host:port")

func (a NetAddress) String() string {
	var res string

	if a.Scheme != "" {
		res += a.Scheme + "://"
	}

	return res + a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, "://")

	if len(hp) == 2 {
		a.Scheme = hp[0]
		s = hp[1]
	}

	hp = strings.Split(s, ":")

	if len(hp) != 2 {
		return errInvalidAddress
	}

	port, err := strconv.Atoi(hp[1])

	if err != nil {
		return fmt.Errorf("failed to convert port '%s' to int: %w", hp[1], err)
	}

	a.Host = hp[0]
	a.Port = port

	return nil
}

func init() {
	flag.Var(&ServerConfig.Server, "a", "http server address")
	flag.Var(&ServerConfig.Base, "b", "base server address")
}

func (c Config) ParseFlags() {
	flag.Parse()

	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		err := c.Server.Set(serverAddress)
		if err != nil {
			fmt.Printf("Failed to set server address '%s' in config", serverAddress)
			fmt.Println(err)
		}
	}

	if baseAddress := os.Getenv("BASE_ADDRESS"); baseAddress != "" {
		err := c.Server.Set(baseAddress)
		if err != nil {
			fmt.Printf("Failed to set base address '%s' in config", baseAddress)
			fmt.Println(err)
		}
	}
}
