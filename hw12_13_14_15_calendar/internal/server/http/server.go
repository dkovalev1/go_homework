package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"           //nolint
	logger "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger" //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage"       //nolint
)

type contextKey string

var keyServerAddr contextKey = "serverAddr"

type Server struct { // TODO
	log       *logger.Logger
	server    http.Server
	cancelCtx context.CancelFunc
}

type Application interface { // TODO
}

type handler struct { // TODO
	log     *logger.Logger
	storage app.Storage
}

func (h *handler) info(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok\n"))
}

func (h *handler) createEvent(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var event storage.Event
	err := decoder.Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = h.storage.CreateEvent(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}

func (h *handler) deleteEvent(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	err := h.storage.DeleteEvent(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}

func (h *handler) updateEvent(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var event storage.Event
	err := decoder.Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = h.storage.UpdateEvent(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}

type TimeRequest struct {
	Time time.Time
}

type eventSelector func(time.Time) ([]storage.Event, error)

func (h *handler) getAllEvents(w http.ResponseWriter, selector eventSelector) {
	now := time.Now()
	events, err := selector(now)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) getEvents(w http.ResponseWriter, req *http.Request) {
	interval := req.PathValue("interval")

	switch interval {
	case "day":
		h.getAllEvents(w, h.storage.GetAllEventsDay)
	case "week":
		h.getAllEvents(w, h.storage.GetAllEventsWeek)
	case "month":
		h.getAllEvents(w, h.storage.GetAllEventsMonth)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unknown interval"))
	}
}

func (h *handler) ServeHTTP(_ http.ResponseWriter, r *http.Request) { // TODO
	h.log.Info(fmt.Sprintf("Serving %s %s", r.Method, r.URL.Path))
}

func NewServer(port int, logger *logger.Logger, storage app.Storage) *Server {
	h := &handler{log: logger, storage: storage}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", loggingMiddleware(logger, h.info))
	mux.HandleFunc("GET /info", loggingMiddleware(logger, h.info))
	mux.HandleFunc("GET /hello", loggingMiddleware(logger, h.info))

	mux.HandleFunc("PUT /event", loggingMiddleware(logger, h.createEvent))
	mux.HandleFunc("POST /event", loggingMiddleware(logger, h.updateEvent))
	mux.HandleFunc("DELETE /event/{id}", loggingMiddleware(logger, h.deleteEvent))
	mux.HandleFunc("GET /event/{interval}", loggingMiddleware(logger, h.getEvents))

	ctx, cancelCtx := context.WithCancel(context.Background())

	return &Server{
		log: logger,
		server: http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           mux,
			ReadHeaderTimeout: 3 * time.Second,
			BaseContext: func(l net.Listener) context.Context {
				ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
				return ctx
			},
		},
		cancelCtx: cancelCtx,
	}
}

func (s *Server) Start(_ context.Context) error {
	s.log.Info(
		fmt.Sprintf("starting REST server on %s", s.server.Addr),
	)

	go func() {
		if err := s.server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			s.log.Info(err.Error())
		}
		s.cancelCtx()
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	return nil
}
