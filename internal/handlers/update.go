package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/artgromov/observer/internal/storage"
	"github.com/go-chi/chi/v5"
)

type UpdateMetricsHandler struct {
	Storage storage.Storage
}

func (mh *UpdateMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "metric_type")
	metricName := chi.URLParam(r, "metric_name")
	metricValueString := chi.URLParam(r, "metric_value")

	logger.Printf("updating metricType: %s, metricName: %s, metricValue: %s", metricType, metricName, metricValueString)

	if metricName == "" {
		http.Error(w, "invalid URL, metric name must be specified", http.StatusNotFound)
		return
	}

	now := time.Now()
	dateHeader := now.Format(http.TimeFormat)

	switch metricType {
	case "gauge":
		metricValue, err := strconv.ParseFloat(metricValueString, 64)
		if err != nil {
			http.Error(w, "invalid metric value, not float64", http.StatusBadRequest)
			return
		}
		err = mh.Storage.UpdateGauge(metricName, metricValue)
		if err != nil {
			http.Error(w, "failed to update Gauge", http.StatusInternalServerError)
			return
		}
	case "counter":
		metricValue, err := strconv.ParseInt(metricValueString, 10, 64)
		if err != nil {
			http.Error(w, "invalid metric value, not int64", http.StatusBadRequest)
			return
		}
		err = mh.Storage.UpdateCounter(metricName, metricValue)
		if err != nil {
			http.Error(w, "failed to update Counter", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid metric type, only gauge or counter are supported", http.StatusBadRequest)
		return
	}
	w.Header().Set("Date", dateHeader)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
