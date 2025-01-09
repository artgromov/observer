package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/artgromov/observer/internal/client"
	"github.com/artgromov/observer/internal/collectors"
	"github.com/artgromov/observer/internal/configs"
	"github.com/artgromov/observer/internal/logger"
	"go.uber.org/zap"
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

	l := logger.Get()
	buildInfo, ok := debug.ReadBuildInfo()
	var gitRevision string
	if ok {
		for _, v := range buildInfo.Settings {
			if v.Key == "vcs.revision" {
				gitRevision = v.Value
				break
			}
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	l.Info("starting agent", zap.String("server_addr", cfg.Addr), zap.String("git_revision", gitRevision), zap.String("go_version", buildInfo.GoVersion))

	cl := client.NewClient(fmt.Sprintf("http://%s", cfg.Addr))

	rcl := collectors.NewRuntimeCollector(cl, time.Duration(cfg.PollInterval)*time.Second, time.Duration(cfg.ReportInterval)*time.Second)

	rcl.Start()
	<-sigs
	rcl.Stop()
}
