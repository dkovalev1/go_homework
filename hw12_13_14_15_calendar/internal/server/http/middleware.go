package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func logRequest(r *http.Request, duration time.Duration) {
	// TODO
	fmt.Printf("%s [%s] %s %s %v %v \"%s\"\n",
		r.RemoteAddr,
		time.Now().Format(time.RFC3339),
		r.Method,
		r.URL.Path,
		r.Proto,
		duration,
		r.UserAgent(),
	)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		start := time.Now()
		next(w, r)

		duration := time.Since(start)

		logRequest(r, duration)
	})
}
