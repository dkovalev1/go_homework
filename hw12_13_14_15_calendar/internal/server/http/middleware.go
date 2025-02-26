package internalhttp

import (
	"fmt"
	"net/http"
	"time"
	logger "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger" //nolint
)

func logRequest(r *http.Request, duration time.Duration, logger *logger.Logger) {
	logger.Info(
		fmt.Sprintf("%s [%s] %s %s %v %v \"%s\"\n",
			r.RemoteAddr,
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			r.Proto,
			duration,
			r.UserAgent(),
		),
	)
}

func loggingMiddleware(logger *logger.Logger, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		start := time.Now()
		next(w, r)

		duration := time.Since(start)

		logRequest(r, duration, logger)
	})
}
