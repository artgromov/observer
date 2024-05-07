package configs

import (
	"os"

	flag "github.com/spf13/pflag"
)

type ServerConfig struct {
	Addr string
}

func NewServerConfig() *ServerConfig {
	c := new(ServerConfig)
	flag.StringVarP(&c.Addr, "addr", "a", "localhost:8080", "addr to use for server")
	return c
}

func (c *ServerConfig) Parse() error {
	flag.Parse()
	value, exists := os.LookupEnv("ADDRESS")
	if exists {
		c.Addr = value
	}
	return validateAddrString(c.Addr)
}
