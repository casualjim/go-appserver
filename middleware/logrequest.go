package middleware

import (
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
)

type Logger interface {
	Printf(string, ...interface{})
	Fatalf(string, ...interface{})
}

func LogRequests(lg Logger) func(http.Handler) http.Handler {
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
				lg.Printf(
					"%s request host=%s proto=%s method=%s path=%s status=%d took=%s",
					r.URL.Scheme,
					r.RemoteAddr,
					r.Proto,
					r.Method,
					r.RequestURI,
					status,
					time.Since(start).String(),
				)
			}()
			next.ServeHTTP(rw, r)
		})
	}
}
