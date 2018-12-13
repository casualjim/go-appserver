package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestRecover(t *testing.T) {
	var buf bytes.Buffer
	handler := Recover(log.New(&buf, "", 0))
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		panic("Unexpected error!")
	})

	recovery := handler(handlerFunc)
	recovery.ServeHTTP(httptest.NewRecorder(), newRequest("GET", "/subdir/asdf"))

	if !strings.Contains(buf.String(), "Unexpected error!") {
		t.Fatalf("Got log %#v, wanted substring %#v", buf.String(), "Unexpected error!")
	}
}
