package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	storage "github.com/artgromov/observer/internal/storage"
)

type UpdateMetricsHandler struct {
	Storage storage.Storage
}

func (mh *UpdateMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	if !strings.HasPrefix(path, "/update/") {
		panic(fmt.Errorf("unexpected URL path prefix: %s", path))
	}

	path, _ = strings.CutPrefix(path, "/update/")
	pathComponents := strings.Split(path, "/")
	logger.Printf("got URL pathComponents %s, len %d", &pathComponents, len(pathComponents))

	if len(pathComponents) != 3 {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	metricType := pathComponents[0]
	metricName := pathComponents[1]
	metricStringValue := pathComponents[2]

	if metricName == "" {
		http.Error(w, "invalid URL, metric name must be specified", http.StatusNotFound)
		return
	}

	now := time.Now()
	dateHeader := now.Format(http.TimeFormat)

	switch metricType {
	case "gauge":
		metricValue, err := strconv.ParseFloat(metricStringValue, 64)
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
		metricValue, err := strconv.ParseInt(metricStringValue, 10, 64)
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
