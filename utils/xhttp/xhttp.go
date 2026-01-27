package xhttp

import (
	"context"
	"fmt"
	"github.com/og-saas/framework/metadata"
	"github.com/og-saas/framework/utils/xerr"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/x/errors"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/status"
)

const (
	BusinessCodeOK = 0
	BusinessMsgOk  = "ok"

	serverTimeHeader = "X-Server-Time"
)

type BaseResponse[T any] struct {
	Code    int    `json:"code" xml:"code"`
	Message string `json:"message" xml:"message"`
	Data    T      `json:"data,omitempty" xml:"data,omitempty"`
	TraceID string `json:"trace_id,omitempty" xml:"trace_id,omitempty"`
}

// JsonBaseResponseCtx writes v into w with appropriate http status code.
func JsonBaseResponseCtx(ctx context.Context, w http.ResponseWriter, v any) {
	w.Header().Set(serverTimeHeader, fmt.Sprintf("%d", time.Now().Unix()))
	httpx.OkJsonCtx(ctx, w, wrapBaseResponse(ctx, v))
}
func wrapBaseResponse(ctx context.Context, v any) BaseResponse[any] {
	var resp BaseResponse[any]
	switch data := v.(type) {
	case xerr.Error:
		resp.Code = data.Code.Int()
		resp.Message = data.GetMessage(metadata.GetLanguageFromCtx(ctx))
		resp.Data = data.Data
	case errors.CodeMsg:
		resp.Code = data.Code
		resp.Message = data.Msg
	case *status.Status:
		resp.Code = int(data.Code())
		resp.Message = data.Message()
	case error:
		resp.Code = http.StatusInternalServerError
		resp.Message = data.Error()
	default:
		resp.Code = BusinessCodeOK
		resp.Message = BusinessMsgOk
		resp.Data = v
	}
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		resp.TraceID = spanCtx.TraceID().String()
	}
	return resp
}
