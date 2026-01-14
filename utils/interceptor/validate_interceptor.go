package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		if v, ok := req.(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		return handler(ctx, req)
	}
}

func ValidateStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		return handler(srv, &wrappedValidateServerStream{ServerStream: ss, ctx: ss.Context()})
	}
}

type wrappedValidateServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (wrapped *wrappedValidateServerStream) RecvMsg(m interface{}) error {
	if err := wrapped.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	if vld, ok := m.(interface{ Validate() error }); ok {
		if err := vld.Validate(); err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
	}

	return nil
}

func (wrapped *wrappedValidateServerStream) Context() context.Context {
	return wrapped.ctx
}
