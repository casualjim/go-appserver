package middleware

import (
	"github.com/stretchr/testify/assert"
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

func TestSetAcceptEncodingForPushOptionsWithoutHeaders(t *testing.T) {
	var opts *http.PushOptions
	opts = setAcceptEncodingForPushOptions(opts)

	assert.NotNil(t, opts)
	assert.NotNil(t, opts.Header)

	for k, v := range opts.Header {
		assert.Equal(t, "Accept-Encoding", k)
		assert.Len(t, v, 1)
		assert.Equal(t, "gzip", v[0])
	}

	opts = &http.PushOptions{}
	opts = setAcceptEncodingForPushOptions(opts)

	assert.NotNil(t, opts)
	assert.NotNil(t, opts.Header)

	for k, v := range opts.Header {
		assert.Equal(t, "Accept-Encoding", k)
		assert.Len(t, v, 1)
		assert.Equal(t, "gzip", v[0])
	}
}


// setAcceptEncodingForPushOptions sets "Accept-Encoding" : "gzip" for PushOptions without overriding existing headers.
func setAcceptEncodingForPushOptions(opts *http.PushOptions) *http.PushOptions {

	if opts == nil {
		opts = &http.PushOptions{
			Header: http.Header{
				"Accept-Encoding": []string{"gzip"},
			},
		}
		return opts
	}

	if opts.Header == nil {
		opts.Header = http.Header{
			"Accept-Encoding": []string{"gzip"},
		}
		return opts
	}

	if encoding := opts.Header.Get("Accept-Encoding"); encoding == "" {
		opts.Header.Add("Accept-Encoding", "gzip")
		return opts
	}

	return opts
}


func TestSetAcceptEncodingForPushOptionsWithHeaders(t *testing.T) {
	opts := &http.PushOptions{
		Header: http.Header{
			"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.98 Safari/537.36"},
		},
	}
	opts = setAcceptEncodingForPushOptions(opts)

	assert.NotNil(t, opts)
	assert.NotNil(t, opts.Header)

	assert.Equal(t, "gzip", opts.Header.Get("Accept-Encoding"))
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.98 Safari/537.36", opts.Header.Get("User-Agent"))

	opts = &http.PushOptions{
		Header: http.Header{
			"User-Agent":   []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.98 Safari/537.36"},
			"Accept-Encoding": []string{"deflate"},
		},
	}
	opts = setAcceptEncodingForPushOptions(opts)

	assert.NotNil(t, opts)
	assert.NotNil(t, opts.Header)

	e, found := opts.Header["Accept-Encoding"]
	if !found {
		assert.Fail(t, "Missing Accept-Encoding header value")
	}
	assert.Len(t, e, 1)
	assert.Equal(t, "deflate", e[0])
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.98 Safari/537.36", opts.Header.Get("User-Agent"))
}
