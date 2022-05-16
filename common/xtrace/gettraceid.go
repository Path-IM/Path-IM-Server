package xtrace

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

func TraceIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		traceId := spanCtx.TraceID().String()
		//logx.WithContext(ctx).Info("获取 trace-id ", traceId)
		return traceId
	}

	return ""
}
