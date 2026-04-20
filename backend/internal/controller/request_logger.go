package controller

import (
	"log"
	"net/http"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (c *Controller) RequestLogger(next http.Handler) http.Handler {
	_ = c

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		writer := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(writer, r)

		log.Printf("http_request method=%s path=%s status=%d duration=%s remote=%s", r.Method, r.URL.Path, writer.statusCode, time.Since(startedAt), r.RemoteAddr)
	})
}
