package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/rutkin/url-shortener/internal/app/logger"
	"go.uber.org/zap"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	contentType := w.Header().Get("Content-Type")
	if contentType == "application/json" || contentType == "text/html" {
		w.Header().Set("Content-Encoding", "gzip")
		return w.Writer.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func WithCompress(h http.Handler) http.Handler {
	compFn := func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gz := gzip.NewWriter(w)
			defer gz.Close()
			ow = gzipWriter{ResponseWriter: w, Writer: gz}
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
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
