package tracing

import (
	"github.com/openzipkin/zipkin-go"
)

func ExampleNewTracer() {
	enpoitUrl := "http://tracing-analysis-dc-qd.aliyuncs.com/adapt_a92srsbtkl@xxx@53df7ad2afe8301/api/v2/spans"
	tracer := GetTracer("manual_demoService", "172.20.23.100:80", enpoitUrl)
	// tracer can now be used to create spans.
	span := tracer.StartSpan("some_operation")
	// ... do some work ...
	span.Finish()

	childSpan := tracer.StartSpan("some_operation2", zipkin.Parent(span.Context()))
	// ... do some work ...
	childSpan.Finish()

	span.Finish()

	// Output:
}
