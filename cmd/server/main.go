package main

import (
	"log"
	"net/http"
	"os"

	"github.com/artgromov/observer/internal/configs"
	"github.com/artgromov/observer/internal/server"
	"github.com/artgromov/observer/internal/storage"
)

var logger = log.New(os.Stdout, "", 0)

func main() {
	cfg := configs.NewServerConfig()
	err := cfg.Parse()
	if err != nil {
		panic(err)
	}

	ms := storage.NewMemStorage()

	r := server.MetricsRouter(ms)

	logger.Printf("starting server on addr: \"%s\"", cfg.Addr)
	err = http.ListenAndServe(cfg.Addr, r)
	if err != nil {
		panic(err)
	}
}
