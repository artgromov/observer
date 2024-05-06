package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artgromov/observer/internal/collectors"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	rcl := collectors.NewRuntimeCollector("http://localhost:8080", 2*time.Second, 10*time.Second)

	rcl.Start()
	<-sigs
	rcl.Stop()
}
