package xtrace

import (
	"context"
	"github.com/showurl/Zero-IM-Server/common/utils"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
	"net/http"
)

func RunWithTrace(
	traceId string,
	f func(ctx context.Context),
	kv ...attribute.KeyValue,
) {
	propagator := otel.GetTextMapPropagator()
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	header := http.Header{}
	if traceId != "" {
		header.Set("x-trace-id", traceId)
	}
	ctx := propagator.Extract(context.Background(), propagation.HeaderCarrier(header))
	spanName := utils.CallerFuncName()
	traceIDFromHex, _ := oteltrace.TraceIDFromHex(traceId)
	ctx = oteltrace.ContextWithSpanContext(ctx, oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
		TraceID: traceIDFromHex,
	}))
	spanCtx, span := tracer.Start(
		ctx,
		spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindConsumer),
		oteltrace.WithAttributes(kv...),
	)
	defer span.End()
	propagator.Inject(spanCtx, propagation.HeaderCarrier(header))
	f(spanCtx)
}
