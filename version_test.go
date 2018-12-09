package appserver


import (
	"encoding/json"
	"github.com/casualjim/go-appserver/mocks"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestServeVersion(t *testing.T) {
	info := VersionInfo{
		BuildDate: "today",
		GitCommit: "commit-here",
		GitState:  "clean",
		Version:   "version-here",
	}
	lg := &mocks.LoggerMock{
		EnabledFunc: func() bool { return true },
		ErrorFunc: func(err error, msg string, keysAndValues ...interface{}) {

		},
		InfoFunc: func(msg string, keysAndValues ...interface{}) {

		},
		VFunc: func(level int) logr.InfoLogger {
			return nil
		},
		WithNameFunc: func(name string) logr.Logger {
			return nil
		},
		WithValuesFunc: func(keysAndValues ...interface{}) logr.Logger {
			return nil
		},
	}
	handler := VersionHandler(lg, info)

	b, _ := json.Marshal(info)
	req := httptest.NewRequest("GET", "/version", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	require.Equal(t, 200, resp.Code)
	require.JSONEq(t, string(b), resp.Body.String())
}
