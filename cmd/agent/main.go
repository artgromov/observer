package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artgromov/observer/internal/client"
	"github.com/artgromov/observer/internal/collectors"
	"github.com/artgromov/observer/internal/configs"
)

type Config struct {
	addr string
}

func main() {
	cfg := configs.NewAgentConfig()
	err := cfg.Parse()
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cl := client.NewClient(fmt.Sprintf("http://%s", cfg.Addr))

	rcl := collectors.NewRuntimeCollector(cl, time.Duration(cfg.PollInterval)*time.Second, time.Duration(cfg.ReportInterval)*time.Second)

	rcl.Start()
	<-sigs
	rcl.Stop()
}
