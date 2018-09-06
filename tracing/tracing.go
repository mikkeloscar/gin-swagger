package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

const (
	spanContextKey = "opentracing_span"
)

// InitSpan initializes a new span. It tries to extract parent span from the
// HTTP headers of the request and will be initialized as a child span if it
// succeeds. Otherwise it will start a new span which is not a child.
func InitSpan(tracer opentracing.Tracer, operationName string, opts ...opentracing.StartSpanOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		spanCtx, err := tracer.Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
		if err == nil {
			// if we got a span context from the headers,
			// initialize the new span as a child
			opts = append([]opentracing.StartSpanOption{opentracing.ChildOf(spanCtx)}, opts...)
		}
		span := tracer.StartSpan(operationName, opts...)
		ctx.Set(spanContextKey, span)
		defer span.Finish()

		ctx.Next()
	}
}

// Context returns the current tracing context for the request. If no span is
// set on the gin.Context it will return a context.Background().
func Context(ctx *gin.Context) context.Context {
	span, ok := ctx.Get(spanContextKey)
	if !ok {
		return context.Background()
	}

	parentSpan, ok := span.(opentracing.Span)
	if !ok {
		return context.Background()
	}

	return opentracing.ContextWithSpan(context.Background(), parentSpan)
}

// StartSpanFromContextWithTracer starts and returns a Span with
// `operationName`, using any Span found within `ctx` as a ChildOfRef. If no
// such parent could be found, StartSpanFromContext creates a root (parentless)
// Span.
//
// The second return value is a context.Context object built around the
// returned Span.
//
// Example usage:
//
//    SomeFunction(ctx context.Context, tracer opentracing.Tracer, ...) {
//        sp, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, "SomeFunction")
//        defer sp.Finish()
//        ...
//    }
func StartSpanFromContextWithTracer(ctx context.Context, tracer opentracing.Tracer, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	var span opentracing.Span
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
		span = tracer.StartSpan(operationName, opts...)
	} else {
		span = tracer.StartSpan(operationName, opts...)
	}
	return span, opentracing.ContextWithSpan(ctx, span)
}
