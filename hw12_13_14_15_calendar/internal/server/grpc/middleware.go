package internalgrpc

import (
	"context"
	"fmt"
	"time"

	logger "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger" //nolint
	"google.golang.org/grpc"                                                         //nolint
)

type CallLogger struct {
	logger *logger.Logger
}

func (logger *CallLogger) logCall(ctx context.Context, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	logger.logger.Info(
		fmt.Sprintf("%s GRPC  Call of %s",
			time.Now().Format(time.RFC3339),
			info.FullMethod))

	resp, err = handler(ctx, req)
	return resp, err
}
