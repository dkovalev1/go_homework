package internalgrpc

import (
	"context"
	"fmt"
	"log"
	"net"

	calendarpb "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/api"           //nolint
	app "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"         //nolint
	logger "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"   //nolint
	storage "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage" //nolint
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CalendarService struct {
	calendarpb.UnimplementedCalendarServer
	storage app.Storage
	lsn     net.Listener
	server  *grpc.Server
	logger  *logger.Logger
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

func (cs *CalendarService) CreateEvent(_ context.Context, req *calendarpb.Event) (*calendarpb.Result, error) {
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

func (cs *CalendarService) UpdateEvent(_ context.Context, req *calendarpb.Event) (*calendarpb.Result, error) {
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

func (cs *CalendarService) DeleteEvent(_ context.Context, req *calendarpb.Event) (*calendarpb.Result, error) {
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

func (cs *CalendarService) GetAllEventsDay(_ context.Context, req *calendarpb.TimeSpec) (*calendarpb.Result, error) {
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

func (cs *CalendarService) GetAllEventsWeek(_ context.Context, req *calendarpb.TimeSpec) (*calendarpb.Result, error) {
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
		Events: makePbEvents(events),
	}, nil
}

func (cs *CalendarService) GetAllEventsMonth(_ context.Context, req *calendarpb.TimeSpec) (*calendarpb.Result, error) {
	events, err := cs.storage.GetAllEventsMonth(req.Stamp.AsTime())
	if err != nil {
		errmsg := err.Error()
		return &calendarpb.Result{
			IsOk:   false,
			Errmsg: &errmsg,
		}, err
	}
	return &calendarpb.Result{
		IsOk:   true,
		Events: makePbEvents(events),
	}, nil
}

// func UnaryServerRequestValidatorInterceptor(logger CallLogger) grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,
// 	) (interface{}, error) {
// 		logger(req)
// 		return handler(ctx, req)
// 	}
// }

func NewService(port int, logger *logger.Logger, storage app.Storage) *CalendarService {
	// For educational project let's relax security requirements for now and bind
	// to all interfaces.
	address := fmt.Sprintf(":%d", port)
	lsn, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	grpcLogger := CallLogger{logger: logger}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpcLogger.logCall),
	)
	ret := &CalendarService{server: server, storage: storage, lsn: lsn, logger: logger}

	calendarpb.RegisterCalendarServer(server, ret)

	return ret
}

func (cs *CalendarService) serve() {
	if err := cs.server.Serve(cs.lsn); err != nil {
		log.Fatal(err)
	}
}

func RunServer(server *CalendarService) {
	server.serve()
}

func (cs *CalendarService) Start() error {
	cs.logger.Info(
		fmt.Sprintf("starting grpc server on %s", cs.lsn.Addr().String()),
	)

	go func() {
		cs.server.Serve(cs.lsn)
	}()

	return nil
}

func (cs *CalendarService) Stop() error {
	cs.server.Stop()
	return nil
}
