package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/rutkin/url-shortener/internal/app/logger"
	"go.uber.org/zap"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w *gzipWriter) useCompression() bool {
	contentType := w.Header().Get("Content-Type")
	return contentType == "application/json" || contentType == "text/html"
}

// write body with gzip comprassion
func (w *gzipWriter) Write(b []byte) (int, error) {
	if w.useCompression() {
		w.ResponseWriter.Header().Add("Content-Encoding", "gzip")
		return w.Writer.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

// write header with gzip
func (w *gzipWriter) WriteHeader(statusCode int) {
	if w.useCompression() {
		w.ResponseWriter.Header().Add("Content-Encoding", "gzip")
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

// close
func (w *gzipWriter) Close() error {
	if w.useCompression() {
		return w.Writer.Close()
	}
	return nil
}

// comprassion middleware
func WithCompress(h http.Handler) http.Handler {
	compFn := func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			logger.Log.Info("Using gzip writer for request")
			gz := &gzipWriter{ResponseWriter: w, Writer: gzip.NewWriter(w)}
			defer gz.Close()
			ow = gz
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			logger.Log.Info("Using gzip reader for request")
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Log.Error("failed to create gzip reader", zap.String("error", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
			}
			r.Body = gr
			defer gr.Close()
		}

		h.ServeHTTP(ow, r)
	}
	return http.HandlerFunc(compFn)
}
