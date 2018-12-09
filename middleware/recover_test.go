package middleware

import (
	"bytes"
	"fmt"
	"github.com/casualjim/go-appserver/mocks"
	"github.com/go-logr/logr"
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

	lg := &mocks.LoggerMock{
		EnabledFunc: func() bool { return true },
		ErrorFunc: func(err error, msg string, keysAndValues ...interface{}) {
			buf.WriteString(fmt.Sprintf("msg[%v]: %v", fmt.Sprint(keysAndValues), err))
		},
		InfoFunc: func(msg string, keysAndValues ...interface{}) {
			buf.WriteString(fmt.Sprintf("msg: %v", fmt.Sprint(keysAndValues)))
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
	handler := Recover(lg)
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		panic("Unexpected error!")
	})

	recovery := handler(handlerFunc)
	recovery.ServeHTTP(httptest.NewRecorder(), newRequest("GET", "/subdir/asdf"))

	if !strings.Contains(buf.String(), "Unexpected error!") {
		t.Fatalf("Got log %#v, wanted substring %#v", buf.String(), "Unexpected error!")
	}
}
