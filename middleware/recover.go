package middleware

import (
	"net/http"
	"runtime/debug"
)

func Recover(lg Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if err, ok := rvr.(error); ok {
						lg.Printf("%v\n%s", err, string(debug.Stack()))
					} else {
						lg.Printf("%v\n%s", rvr, string(debug.Stack()))
					}
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
