package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/rybalka1/devmetrics/internal/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateRoutes(t *testing.T) {
	var (
		store = memstorage.NewMemStorage()
		route = NewRouter(store)
		srv   = httptest.NewServer(route)
	)

	updateURL := "/update"
	tests := []struct {
		name       string
		store      memstorage.Storage
		haveURL    string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "Wrong url 1",
			haveURL:    "/aaaaa",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "Wrong url 2",
			haveURL:    updateURL + "/aaaaa",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "Wrong url 3",
			haveURL:    updateURL + "/counter/sys",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "Wrong url 4",
			haveURL:    updateURL + "/counter/sys/",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "Wrong Type",
			haveURL:    updateURL + "/aaaaa/sys/123",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "Wrong Val Counter",
			haveURL:    updateURL + "/counter/sys/213.214",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "Good Type Counter",
			haveURL:    updateURL + "/counter/sys/213",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Good Type Gauge",
			haveURL:    updateURL + "/gauge/sys/213.11",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Good Type Gauge 2",
			haveURL:    updateURL + "/gauge/sys/213",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status, _, err := testSendRequest(srv, "POST", test.haveURL, nil)
			require.NoError(t, err)

			assert.Equal(t, test.statusCode, status)

			if test.wantErr == false {
				filterURL, _ := strings.CutPrefix(test.haveURL, updateURL)
				filterURL = strings.Trim(filterURL, "/")
				params := strings.Split(filterURL, "/")

				if len(params) != 3 {
					t.Error("wrong parsing url")
					return
				}

				mType := params[0]
				mName := params[1]
				expectValue := params[2]
				actualValue := store.GetMetricString(mType, mName)
				fmt.Println(filterURL)
				assert.Equal(t, expectValue, actualValue)

			}

		})
	}
}

func testSendRequest(server *httptest.Server, method, url string,
	data []byte) (respCode int, body []byte, err error) {
	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest(method, server.URL+url, buf)
	if err != nil {
		return 0, nil, err
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil

}

// type HTTPparams struct {
// 	method string
// 	url    string
// 	data   string
// }

func TestGetMetric(t *testing.T) {
	var (
		store = memstorage.NewMemStorage()
		route = NewRouter(store)
		srv   = httptest.NewServer(route)
	)
	testgauge := 121.99
	testcounter := 122

	strtestgauge := strconv.FormatFloat(testgauge, 'f', -1, 64)
	store.UpdateGauges("testgauge", testgauge)
	strtestcounter := strconv.Itoa(testcounter)
	store.UpdateCounters("testcounter", int64(testcounter))

	defaultURL := "/value"

	tests := []struct {
		name         string
		url          string
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "No metric name",
			url:          defaultURL + "/counter/",
			wantStatus:   http.StatusNotFound,
			wantResponse: NotFound,
		},
		{
			name:         "Wrong metric name",
			url:          defaultURL + "/counter/aaaa",
			wantStatus:   http.StatusNotFound,
			wantResponse: NotFound,
		},
		{
			name:         "Good counter",
			url:          defaultURL + "/counter/testcounter",
			wantStatus:   http.StatusOK,
			wantResponse: strtestcounter,
		},
		{
			name:         "Good gauge",
			url:          defaultURL + "/gauge/testgauge",
			wantStatus:   http.StatusOK,
			wantResponse: strtestgauge,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status, body, err := testSendRequest(srv, "GET", test.url, nil)
			require.NoError(t, err)
			assert.Equal(t, test.wantStatus, status)
			assert.Equal(t, test.wantResponse, string(body))
		})
	}
}
