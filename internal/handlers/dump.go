package handlers

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/artgromov/observer/internal/storage"
)

type DumpMetricsHandler struct {
	Storage storage.Storage
}

func (mh *DumpMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	dateHeader := now.Format(http.TimeFormat)

	result := mh.Storage.Dump()

	w.Header().Set("Date", dateHeader)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, strings.Join(result, "\n"))
}
