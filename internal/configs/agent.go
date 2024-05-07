package configs

import (
	flag "github.com/spf13/pflag"
)

type AgentConfig struct {
	Addr           string
	ReportInterval uint64
	PollInterval   uint64
}

func NewAgentConfig() *AgentConfig {
	c := new(AgentConfig)
	flag.StringVarP(&c.Addr, "addr", "a", "localhost:8080", "addr to use for server")
	flag.Uint64VarP(&c.ReportInterval, "report-interval", "r", 10, "runtime metrics report interval")
	flag.Uint64VarP(&c.PollInterval, "poll-interval", "p", 2, "runtime metrics poll interval")
	return c
}

func (c *AgentConfig) Parse() error {
	flag.Parse()
	return validateAddrString(c.Addr)
}
