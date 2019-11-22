package weather

import (
	"net/http"

	"context"

	"github.com/knative-sample/cloud-native-go-weather/pkg/tracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

func (wa *WebApi) TraceLog(ha http.Handler) http.Handler {
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		wa.tracer,
		zipkinhttp.TagResponseSize(true),
		zipkinhttp.SpanName("weather"),
	)
	return serverMiddleware(ha)
}

func (wa *WebApi) NewSpan(name string, ctx context.Context) zipkin.Span {
	var currentSpan zipkin.Span
	tracer := tracing.GetTracer(wa.ServiceName, wa.InstanceIp, wa.ZipKinEndpoint)
	if parent := zipkin.SpanFromContext(ctx); parent != nil {
		currentSpan = tracer.StartSpan(name, zipkin.Parent(parent.Context()))
	} else {
		currentSpan = tracer.StartSpan(name)
	}
	return currentSpan
}
