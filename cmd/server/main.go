package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger = log.New(os.Stderr, "", 0)

type Storage interface {
	GetGaugeMap() map[string]float64
	UpdateGauge(name string, value float64) error
	GetCounterMap() map[string]int64
	UpdateCounter(name string, value int64) error
}

type MemStorage struct {
	gaugeMap   map[string]float64
	counterMap map[string]int64
}

func NewMemStorage() *MemStorage {
	s := new(MemStorage)
	s.gaugeMap = make(map[string]float64)
	s.counterMap = make(map[string]int64)
	return s
}

func (s *MemStorage) GetGaugeMap() map[string]float64 {
	return s.gaugeMap
}

func (s *MemStorage) UpdateGauge(name string, value float64) error {
	s.gaugeMap[name] = value
	return nil
}

func (s *MemStorage) GetCounterMap() map[string]int64 {
	return s.counterMap
}

func (s *MemStorage) UpdateCounter(name string, value int64) error {
	_, ok := s.counterMap[name]
	if ok {
		s.counterMap[name] += value
	} else {
		s.counterMap[name] = value
	}
	return nil
}

type UpdateMetricsHandler struct {
	storage Storage
}

func (mh *UpdateMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-type") != "text/plain" {
		http.Error(w, "invalid Content-type, text/plain is expected", http.StatusBadRequest)
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
		err = mh.storage.UpdateGauge(metricName, metricValue)
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
		err = mh.storage.UpdateCounter(metricName, metricValue)
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

type GetMetricsHandler struct {
	storage Storage
}

func (mh *GetMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	gm, err := json.Marshal(mh.storage.GetGaugeMap())
	if err != nil {
		http.Error(w, "failed to marshal gauge map to json", http.StatusInternalServerError)
		return
	}
	cm, err := json.Marshal(mh.storage.GetCounterMap())
	if err != nil {
		http.Error(w, "failed to marshal counter map to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(gm)
	_, _ = w.Write(cm)
}

func main() {
	memStorage := NewMemStorage()

	umh := UpdateMetricsHandler{memStorage}
	gmh := GetMetricsHandler{memStorage}

	mux := http.NewServeMux()
	mux.Handle(`/update/`, &umh)
	mux.Handle(`/get/`, &gmh)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
