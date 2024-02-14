package config

import (
	"flag"
	"net/url"
)

type Config struct {
	Address url.URL
}

var ServerConfig = Config{url.URL{Scheme: "http", Host: "localhost:8080"}}

func (c Config) String() string {
	return c.Address.String()
}

func (c *Config) Set(s string) error {
	url, err := url.Parse(s)
	if err != nil {
		return err
	}
	c.Address = *url
	return nil
}

func init() {
	flag.Var(&ServerConfig, "a", "http server address")
	flag.Var(&ServerConfig, "b", "base server address")
}
