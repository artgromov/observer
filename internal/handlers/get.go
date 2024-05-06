package handlers

import (
	"encoding/json"
	"net/http"

	storage "github.com/artgromov/observer/internal/storage"
)

type GetMetricsHandler struct {
	Storage storage.Storage
}

func (mh *GetMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	gm, err := json.Marshal(mh.Storage.GetGaugeMap())
	if err != nil {
		http.Error(w, "failed to marshal gauge map to json", http.StatusInternalServerError)
		return
	}
	cm, err := json.Marshal(mh.Storage.GetCounterMap())
	if err != nil {
		http.Error(w, "failed to marshal counter map to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(gm)
	_, _ = w.Write(cm)
}
