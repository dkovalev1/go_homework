package internalgrpc

import (
	"context"
	"log"
	"net"

	calendarpb "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/api"
	app "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"
	logger "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger" //nolint
	storage "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CalendarService struct {
	calendarpb.UnimplementedCalendarServer
	storage app.Storage
	lsn     net.Listener
	server  *grpc.Server
}

func makeEvent(req *calendarpb.Event) storage.Event {
	return storage.Event{
		ID:          req.ID,
		Title:       req.Title,
		StartTime:   req.StartTime.AsTime(),
		Duration:    req.Duration.AsDuration(),
		Description: req.Description,
		UserID:      req.UserID,
		NotifyTime:  req.NotifyTime.AsDuration(),
	}
}

func makePbEvents(events []storage.Event) []*calendarpb.Event {
	out := make([]*calendarpb.Event, 0)
	for _, event := range events {
		out = append(out, &calendarpb.Event{
			ID:          event.ID,
			Title:       event.Title,
			StartTime:   timestamppb.New(event.StartTime),
			Duration:    durationpb.New(event.Duration),
			Description: event.Description,
			UserID:      event.UserID,
			NotifyTime:  durationpb.New(event.NotifyTime),
		})
	}

	return out
}

func (cs *CalendarService) CreateEvent(ctx context.Context, req *calendarpb.Event) (*calendarpb.Result, error) {
	err := cs.storage.CreateEvent(makeEvent(req))

	if err != nil {
		errmsg := err.Error()
		return &calendarpb.Result{
			IsOk:   false,
			Errmsg: &errmsg,
		}, err
	}
	return &calendarpb.Result{
		IsOk: true,
	}, nil
}

func (cs *CalendarService) UpdateEvent(ctx context.Context, req *calendarpb.Event) (*calendarpb.Result, error) {
	err := cs.storage.UpdateEvent(makeEvent(req))

	if err != nil {
		errmsg := err.Error()
		return &calendarpb.Result{
			IsOk:   false,
			Errmsg: &errmsg,
		}, err
	}
	return &calendarpb.Result{IsOk: true}, nil
}

func (cs *CalendarService) DeleteEvent(ctx context.Context, req *calendarpb.Event) (*calendarpb.Result, error) {
	err := cs.storage.DeleteEvent(makeEvent(req))

	if err != nil {
		errmsg := err.Error()
		return &calendarpb.Result{
			IsOk:   false,
			Errmsg: &errmsg,
		}, err
	}
	return &calendarpb.Result{IsOk: true}, nil
}

func (cs *CalendarService) GetAllEventsDay(ctx context.Context, req *calendarpb.TimeSpec) (*calendarpb.Result, error) {
	events, err := cs.storage.GetAllEventsDay(req.Stamp.AsTime())

	if err != nil {
		errmsg := err.Error()
		return &calendarpb.Result{
			IsOk:   false,
			Errmsg: &errmsg,
		}, err
	}

	return &calendarpb.Result{IsOk: true, Events: makePbEvents(events)}, nil
}

func (cs *CalendarService) GetAllEventsWeek(ctx context.Context, req *calendarpb.TimeSpec) (*calendarpb.Result, error) {
	events, err := cs.storage.GetAllEventsWeek(req.Stamp.AsTime())

	if err != nil {
		errmsg := err.Error()
		return &calendarpb.Result{
			IsOk:   false,
			Errmsg: &errmsg,
		}, err
	}

	return &calendarpb.Result{
		IsOk:   true,
		Events: makePbEvents(events)}, nil
}

func (cs *CalendarService) GetAllEventsMonth(ctx context.Context, req *calendarpb.TimeSpec) (*calendarpb.Result, error) {
	events, err := cs.storage.GetAllEventsMonth(req.Stamp.AsTime())

	if err != nil {
		errmsg := err.Error()
		return &calendarpb.Result{
			IsOk:   false,
			Errmsg: &errmsg,
		}, err
	}
	return &calendarpb.Result{IsOk: true,
		Events: makePbEvents(events)}, nil
}

type CallLogger func(req interface{}) error

func UnaryServerRequestValidatorInterceptor(logger CallLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger(req)
		return handler(ctx, req)
	}
}

func NewService(logger *logger.Logger, storage app.Storage) *CalendarService {

	lsn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(callLogger),
	)
	ret := &CalendarService{server: server, storage: storage, lsn: lsn}

	calendarpb.RegisterCalendarServer(server, ret)

	return ret
}

func (s *CalendarService) serve() {
	if err := s.server.Serve(s.lsn); err != nil {
		log.Fatal(err)
	}
}

func RunServer(server *CalendarService) {
	server.serve()
}

func (s *CalendarService) Start() error {
	log.Printf("starting grpc server on %s", s.lsn.Addr().String())

	go func() {
		s.server.Serve(s.lsn)
	}()

	return nil
}

func (s *CalendarService) Stop() error {
	// TODO
	s.server.Stop()
	return nil
}
