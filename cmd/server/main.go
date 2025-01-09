package main

import (
	"net/http"
	"runtime/debug"

	"github.com/artgromov/observer/internal/configs"
	"github.com/artgromov/observer/internal/logger"
	"github.com/artgromov/observer/internal/server"
	"github.com/artgromov/observer/internal/storage"
	"go.uber.org/zap"
)

func main() {
	cfg := configs.NewServerConfig()
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

	ms := storage.NewMemStorage()

	r := server.MetricsRouter(ms)

	l.Info("starting server", zap.String("server_addr", cfg.Addr), zap.String("git_revision", gitRevision), zap.String("go_version", buildInfo.GoVersion))
	err = http.ListenAndServe(cfg.Addr, r)
	if err != nil {
		panic(err)
	}
}
