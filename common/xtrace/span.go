package xtrace

import (
	"context"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func StartFuncSpan(ctx context.Context, name string, f func(context.Context), kv ...attribute.KeyValue) {
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	spanCtx, span := tracer.Start(ctx, name,
		oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
		oteltrace.WithAttributes(kv...),
	)
	defer span.End()
	f(spanCtx)
}
