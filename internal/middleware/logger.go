package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			defer func() {
				if rec := recover(); rec != nil {
					logger.WithFields(logrus.Fields{
						"method": r.Method,
						"url":    r.URL.Path,
						"panic":  rec,
					}).Error("Panic en request")
				}
			}()

			next.ServeHTTP(w, r)

			logger.WithFields(logrus.Fields{
				"method":   r.Method,
				"url":      r.URL.Path,
				"duration": time.Since(start),
			}).Debug("Request completado")
		})
	}
}

// LogStaticFiles logs requests to static files for debugging
func LogStaticFiles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"path": r.URL.Path,
		}).Debug("Serving static file")
		next.ServeHTTP(w, r)
	})
}
