package xtrace

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/common/utils"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func RunWithTrace(
	traceId string,
	f func(ctx context.Context),
	kv ...attribute.KeyValue,
) {
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	spanName := utils.CallerFuncName()
	traceIDFromHex, _ := oteltrace.TraceIDFromHex(traceId)
	ctx := oteltrace.ContextWithSpanContext(context.Background(), oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
		TraceID: traceIDFromHex,
	}))
	spanCtx, span := tracer.Start(
		ctx,
		spanName,
		//oteltrace.WithSpanKind(oteltrace.SpanKindConsumer),
		//oteltrace.WithAttributes(kv...),
	)
	defer span.End()
	f(spanCtx)
}
