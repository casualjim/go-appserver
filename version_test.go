package appserver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServeVersion(t *testing.T) {
	info := VersionInfo{
		BuildDate: "today",
		GitCommit: "commit-here",
		GitState:  "clean",
		Version:   "version-here",
	}

	handler := VersionHandler(log.New(ioutil.Discard, "", 0), info)

	b, _ := json.Marshal(info)
	req := httptest.NewRequest("GET", "/version", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	require.Equal(t, 200, resp.Code)
	require.JSONEq(t, string(b), resp.Body.String())
}
