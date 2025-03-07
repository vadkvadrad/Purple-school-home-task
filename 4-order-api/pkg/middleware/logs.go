package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	logrus.SetFormatter(&logrus.TextFormatter{})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		logger := logrus.Fields{
			"status code": wrapper.StatusCode,
			"method":      r.Method,
			"path":        r.URL.Path,
			"time":        time.Since(start),
		}

		if wrapper.StatusCode >= 500 {
			logrus.WithFields(logger).Warn("Something went wrong")
		} else {
			logrus.WithFields(logger).Info("ok")
		}
	})
}
