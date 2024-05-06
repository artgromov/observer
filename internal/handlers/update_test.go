package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	storage "github.com/artgromov/observer/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		url    string
		method string
		want   want
	}{
		{
			name:   "wrong counter method",
			url:    "/update/counter/myCounter/3",
			method: http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "wrong counter url",
			url:    "/update/counter/",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "wrong counter url components count",
			url:    "/update/counter/myCounter/3/4/5",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "wrong counter metrics name",
			url:    "/update/counter//3",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "invalid counter value",
			url:    "/update/counter/myCounter/abc",
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "positive counter",
			url:    "/update/counter/myCounter/3",
			method: http.MethodPost,
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "wrong gauge method",
			url:    "/update/gauge/myGauge/3",
			method: http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "wrong gauge url",
			url:    "/update/gauge/",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "wrong gauge url components count",
			url:    "/update/gauge/myGauge/3/4/5",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "wrong gauge metrics name",
			url:    "/update/gauge//3",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "invalid gauge value",
			url:    "/update/gauge/myGauge/abc",
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "positive gauge",
			url:    "/update/gauge/myGauge/3",
			method: http.MethodPost,
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	ms := storage.NewMemStorage()

	umh := UpdateMetricsHandler{Storage: ms}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()
			umh.ServeHTTP(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			if test.want.response != "" {
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.JSONEq(t, test.want.response, string(resBody))
			}
		})
	}
}
