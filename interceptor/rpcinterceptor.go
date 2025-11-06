package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

func RpcInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		//causeErr := errors.Cause(err)
		
	}
	return
}
