package middleware

import (
	"github.com/felixge/httpsnoop"
	"github.com/go-logr/logr"
	"net/http"
	"time"
)

func LogRequests(lg logr.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()

			status := http.StatusOK
			rw = httpsnoop.Wrap(rw, httpsnoop.Hooks{
				WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						status = code
						next(code)
					}
				},
			})

			defer func() {
				lg.Info(
					"http request",
					"host", r.RemoteAddr,
					"proto", r.Proto,
					"method", r.Method,
					"path", r.RequestURI,
					"status", status,
					"took", time.Since(start).String(),
				)
			}()
			next.ServeHTTP(rw, r)
		})
	}
}
