package main

import (
	"net/http"

	"github.com/artgromov/observer/internal/server"
	"github.com/artgromov/observer/internal/storage"
)

func main() {
	ms := storage.NewMemStorage()

	r := server.MetricsRouter(ms)

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
}
