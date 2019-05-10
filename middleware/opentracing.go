package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/CharlesBases/common/log"
)

type opentracer struct {
}

func NegroniOpentracer() *opentracer {
	return new(opentracer)
}

func (l *opentracer) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx, span, err := traceIntoContext(r.Context(), opentracing.GlobalTracer(), r)
	if err != nil {
		log.Error(err)
	}
	defer span.Finish()
	request := r.WithContext(ctx)
	next(rw, request)
}

func traceIntoContext(ctx context.Context, tracer opentracing.Tracer, r *http.Request) (context.Context, opentracing.Span, error) {
	md := make(map[string]string)
	for k, v := range r.Header {
		md[k] = strings.Join(v, ",")
	}
	name := r.URL.Path

	var sp opentracing.Span
	wireContext, err := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	if err != nil {
		sp = tracer.StartSpan(name)
	} else {
		sp = tracer.StartSpan(name, opentracing.ChildOf(wireContext))
	}
	if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md)); err != nil {
		return nil, nil, err
	}
	ctx = opentracing.ContextWithSpan(ctx, sp)
	// ctx = metadata.NewContext(ctx, md)
	return ctx, sp, nil
}
