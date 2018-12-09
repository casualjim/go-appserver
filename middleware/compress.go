package middleware

import (
	"compress/flate"
	"compress/gzip"
	"github.com/felixge/httpsnoop"
	"net/http"
	"strings"
)

// Adapted from http://github.com/gorilla/handlers
// Their middleware is greedy when it comes to implementing response writer methods
// this version uses httpsnoop to avoid expanding the interface and potentially break libraries
// that make assumptions based on those implementations

// CompressHandler gzip compresses HTTP responses for clients that support it
// via the 'Accept-Encoding' header.
//
// Compressing TLS traffic may leak the page contents to an attacker if the
// page contains user input: http://security.stackexchange.com/a/102015/12208
func CompressHandler(h http.Handler) http.Handler {
	return CompressHandlerLevel(h, gzip.DefaultCompression)
}

// CompressHandlerLevel gzip compresses HTTP responses with specified compression level
// for clients that support it via the 'Accept-Encoding' header.
//
// The compression level should be gzip.DefaultCompression, gzip.NoCompression,
// or any integer value between gzip.BestSpeed and gzip.BestCompression inclusive.
// gzip.DefaultCompression is used in case of invalid compression level.
func CompressHandlerLevel(h http.Handler, level int) http.Handler {
	if level < gzip.DefaultCompression || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	L:
		for _, enc := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
			switch strings.TrimSpace(enc) {
			case "gzip":
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Add("Vary", "Accept-Encoding")

				gw, _ := gzip.NewWriterLevel(w, level)
				defer gw.Close()


				w = httpsnoop.Wrap(w, httpsnoop.Hooks{
					WriteHeader: func(headerFunc httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
						return func(code int) {
							w.Header().Del("Content-Length")
							headerFunc(code)
						}
					},
					Write: func(writeFunc httpsnoop.WriteFunc) httpsnoop.WriteFunc {
						return func(b []byte) (i int, e error) {
							h := w.Header()
							if h.Get("Content-Type") == "" {
								h.Set("Content-Type", http.DetectContentType(b))
							}
							h.Del("Content-Length")

							return gw.Write(b)
						}
					},
				})

				break L
			case "deflate":
				w.Header().Set("Content-Encoding", "deflate")
				w.Header().Add("Vary", "Accept-Encoding")

				fw, _ := flate.NewWriter(w, level)
				defer fw.Close()

				w = httpsnoop.Wrap(w, httpsnoop.Hooks{
					WriteHeader: func(headerFunc httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
						return func(code int) {
							w.Header().Del("Content-Length")
							headerFunc(code)
						}
					},
					Write: func(writeFunc httpsnoop.WriteFunc) httpsnoop.WriteFunc {
						return func(b []byte) (i int, e error) {
							h := w.Header()
							if h.Get("Content-Type") == "" {
								h.Set("Content-Type", http.DetectContentType(b))
							}
							h.Del("Content-Length")

							return fw.Write(b)
						}
					},
				})

				break L
			}
		}

		h.ServeHTTP(w, r)
	})
}
