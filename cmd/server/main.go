package main

import (
	"net/http"

	handlers "github.com/artgromov/observer/internal/handlers"
	storage "github.com/artgromov/observer/internal/storage"
)

func main() {
	ms := storage.NewMemStorage()

	umh := handlers.UpdateMetricsHandler{Storage: ms}
	gmh := handlers.GetMetricsHandler{Storage: ms}

	mux := http.NewServeMux()
	mux.Handle(`/update/`, &umh)
	mux.Handle(`/get/`, &gmh)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
