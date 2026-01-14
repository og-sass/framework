package interceptor

import (
	"context"
	"github.com/og-saas/framework/utils/metadatakey"
	"github.com/og-saas/framework/utils/tenant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

func TenantUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		values := md.Get(metadatakey.TenantIdMetadataKey)
		if len(values) == 0 {
			return handler(ctx, req)
		}

		tenantId, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid tenant-id")
		}

		return handler(tenant.SetTenantId(ctx, tenantId), req)
	}
}

func TenantStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return handler(srv, ss)
		}

		values := md.Get(metadatakey.TenantIdMetadataKey)
		if len(values) == 0 {
			return handler(srv, ss)
		}

		tenantId, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid tenant-id")
		}

		wrapped := &tenantServerStream{
			ServerStream: ss,
			ctx:          tenant.SetTenantId(ss.Context(), tenantId),
		}

		return handler(srv, wrapped)
	}
}

type tenantServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (t *tenantServerStream) Context() context.Context {
	return t.ctx
}
