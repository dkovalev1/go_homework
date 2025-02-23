package internalgrpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func callLogger(ctx context.Context, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	log.Printf("%s GRPC  Call of %s",
		time.Now().Format(time.RFC3339),
		info.FullMethod)

	resp, err = handler(ctx, req)
	return resp, err
}
