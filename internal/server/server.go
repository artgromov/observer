package server

import (
	"net/http"

	handlers "github.com/artgromov/observer/internal/handlers"
	middleware "github.com/artgromov/observer/internal/middleware"
	storage "github.com/artgromov/observer/internal/storage"
	"github.com/go-chi/chi/v5"
)

func MetricsRouter(ms storage.Storage) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestLogMiddleware)
	r.Use(middleware.ResponseLogMiddleware)

	r.Method(http.MethodGet, "/", &handlers.DumpMetricsHandler{Storage: ms})
	r.Route("/value", func(r chi.Router) {
		r.Method(http.MethodGet, "/{metric_type}/{metric_name}", &handlers.ValueMetricsHandler{Storage: ms})
	})
	r.Route("/update", func(r chi.Router) {
		r.Method(http.MethodPost, "/{metric_type}/{metric_name}/{metric_value}", &handlers.UpdateMetricsHandler{Storage: ms})
	})

	return r
}
