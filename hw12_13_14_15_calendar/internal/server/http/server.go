package internalhttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	logger "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger" //nolint
)

type Server struct { // TODO
	log    *logger.Logger
	server http.Server
}

type Application interface { // TODO
}

type handler struct { // TODO
	log *logger.Logger
}

func (h *handler) info(w http.ResponseWriter, r *http.Request) { // TODO
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok\n"))
}

func (h *handler) ServeHTTP(_ http.ResponseWriter, r *http.Request) { // TODO
	h.log.Info(fmt.Sprintf("Serving %s %s", r.Method, r.URL.Path))
}

func NewServer(logger *logger.Logger, _ Application) *Server {
	h := &handler{log: logger}
	mux := http.NewServeMux()
	mux.HandleFunc("/info", loggingMiddleware(h.info))

	return &Server{
		log: logger,
		server: http.Server{
			Addr:              ":8080",
			Handler:           mux,
			ReadHeaderTimeout: 3 * time.Second,
		},
	}
}

func (s *Server) Start(_ context.Context) error {
	// TODO
	log.Fatal(s.server.ListenAndServe())
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	s.server.Shutdown(ctx)
	return nil
}

// TODO
