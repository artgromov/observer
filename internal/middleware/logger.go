package middleware

import (
	"net/http"
	"time"

	"github.com/artgromov/observer/internal/logger"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func RequestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTS := time.Now()
		next.ServeHTTP(w, r)
		logger.Get().Info(
			"request",
			zap.String("uri", r.URL.String()),
			zap.String("method", r.Method),
			zap.Int64("elapsed_ns", int64(time.Since(startTS))),
		)
	})
}

func ResponseLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		logger.Get().Info(
			"response",
			zap.Int("status_code", ww.Status()),
			zap.Int("response_bytes", ww.BytesWritten()),
		)
	})
}
