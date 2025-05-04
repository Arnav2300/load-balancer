package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			originalURL := r.URL.String()
			next.ServeHTTP(lrw, r)
			duration := time.Since(start)

			logger.Info("request details",
				zap.Int("status", lrw.statusCode),
				zap.String("method", r.Method),
				zap.String("url", originalURL),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Duration("duration", duration),
			)
		})
	}
}
