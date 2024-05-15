package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artgromov/observer/internal/server"
	"github.com/artgromov/observer/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp, string(data)
}

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
	ts := httptest.NewServer(server.MetricsRouter(storage.NewMemStorage()))
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, data := testRequest(t, ts, test.method, test.url)
			assert.Equal(t, test.want.code, resp.StatusCode)
			if test.want.response != "" {
				assert.JSONEq(t, test.want.response, data)
			}
			resp.Body.Close() // lolwut? fix statictest
		})
	}
}
