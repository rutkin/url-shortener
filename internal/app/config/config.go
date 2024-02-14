package config

import (
	"errors"
	"flag"
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

func (a NetAddress) String() string {
	var res string
	if a.Scheme != "" {
		res += a.Scheme + "://"
	}
	return res + a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, "://")
	if len(hp) > 1 {
		a.Scheme = hp[0]
		s = hp[1]
	}
	hp = strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
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
		c.Server.Set(serverAddress)
	}

	if baseAddress := os.Getenv("BASE_ADDRESS"); baseAddress != "" {
		c.Server.Set(baseAddress)
	}
}
