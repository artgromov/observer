package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/artgromov/observer/internal/storage"
	"github.com/go-chi/chi/v5"
)

type ValueMetricsHandler struct {
	Storage storage.Storage
}

func (mh *ValueMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "metric_type")
	metricName := chi.URLParam(r, "metric_name")
	var metricValueString string

	now := time.Now()
	dateHeader := now.Format(http.TimeFormat)

	switch metricType {
	case "gauge":
		metricValue, err := mh.Storage.GetGauge(metricName)
		if err != nil {
			http.Error(w, fmt.Sprintf("gauge metric with name \"%s\" not found", metricName), http.StatusNotFound)
			return
		}
		metricValueString = strconv.FormatFloat(metricValue, 'f', -1, 64)
	case "counter":
		metricValue, err := mh.Storage.GetCounter(metricName)
		if err != nil {
			http.Error(w, fmt.Sprintf("counter metric with name \"%s\" not found", metricName), http.StatusNotFound)
			return
		}
		metricValueString = strconv.FormatInt(metricValue, 10)
	default:
		http.Error(w, "invalid metric type, only gauge or counter are supported", http.StatusBadRequest)
		return
	}

	w.Header().Set("Date", dateHeader)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, metricValueString)
}
