package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/klauspost/compress/gzip"
)

// Taken from http://github.com/gorilla/handlers

var contentType = "text/plain; charset=utf-8"

func compressedRequest(w *httptest.ResponseRecorder, compression string) {
	CompressHandlerLevel(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(9*1024))
		w.Header().Set("Content-Type", contentType)
		for i := 0; i < 1024; i++ {
			io.WriteString(w, "Gorilla!\n")
		}
	}), gzip.BestCompression).ServeHTTP(w, &http.Request{
		Method: "GET",
		Header: http.Header{
			"Accept-Encoding": []string{compression},
		},
	})

}

func TestCompressHandlerNoCompression(t *testing.T) {
	w := httptest.NewRecorder()
	compressedRequest(w, "")
	if enc := w.HeaderMap.Get("Content-Encoding"); enc != "" {
		t.Errorf("wrong content encoding, got %q want %q", enc, "")
	}
	if ct := w.HeaderMap.Get("Content-Type"); ct != contentType {
		t.Errorf("wrong content type, got %q want %q", ct, contentType)
	}
	if w.Body.Len() != 1024*9 {
		t.Errorf("wrong len, got %d want %d", w.Body.Len(), 1024*9)
	}
	if l := w.HeaderMap.Get("Content-Length"); l != "9216" {
		t.Errorf("wrong content-length. got %q expected %d", l, 1024*9)
	}
}

func TestCompressHandlerGzip(t *testing.T) {
	w := httptest.NewRecorder()
	compressedRequest(w, "gzip")
	if w.HeaderMap.Get("Content-Encoding") != "gzip" {
		t.Errorf("wrong content encoding, got %q want %q", w.HeaderMap.Get("Content-Encoding"), "gzip")
	}
	if w.HeaderMap.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("wrong content type, got %s want %s", w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
	}
	if w.Body.Len() != 72 {
		t.Errorf("wrong len, got %d want %d", w.Body.Len(), 72)
	}
	if l := w.HeaderMap.Get("Content-Length"); l != "" {
		t.Errorf("wrong content-length. got %q expected %q", l, "")
	}
}

func TestCompressHandlerDeflate(t *testing.T) {
	w := httptest.NewRecorder()
	compressedRequest(w, "deflate")
	if w.HeaderMap.Get("Content-Encoding") != "deflate" {
		t.Fatalf("wrong content encoding, got %q want %q", w.HeaderMap.Get("Content-Encoding"), "deflate")
	}
	if w.HeaderMap.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("wrong content type, got %s want %s", w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
	}
	if w.Body.Len() != 54 {
		t.Fatalf("wrong len, got %d want %d", w.Body.Len(), 54)
	}
}

func TestCompressHandlerGzipDeflate(t *testing.T) {
	w := httptest.NewRecorder()
	compressedRequest(w, "gzip, deflate ")
	if w.HeaderMap.Get("Content-Encoding") != "gzip" {
		t.Fatalf("wrong content encoding, got %q want %q", w.HeaderMap.Get("Content-Encoding"), "gzip")
	}
	if w.HeaderMap.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("wrong content type, got %s want %s", w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
	}
}
