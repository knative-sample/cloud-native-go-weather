package tracing

import (
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gorilla/mux"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

func HttpExample() {
	enpoitUrl := "http://tracing-analysis-dc-qd.aliyuncs.com/xxxx/api/v2/spans"
	tracer := GetTracer("http_demoService", "127.0.01", enpoitUrl)

	// create global zipkin http server middleware
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer, zipkinhttp.TagResponseSize(true),
	)

	// create global zipkin traced http client
	client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	// initialize router
	router := mux.NewRouter()

	// start web service with zipkin http server middleware
	ts := httptest.NewServer(serverMiddleware(router))
	defer ts.Close()

	// set-up handlers
	router.Methods("GET").Path("/some_function").HandlerFunc(someFunc(client, ts.URL))
	router.Methods("POST").Path("/other_function").HandlerFunc(otherFunc(client))

	// initiate a call to some_func
	req, err := http.NewRequest("GET", ts.URL+"/some_function", nil)
	if err != nil {
		log.Fatalf("unable to create http request: %+v\n", err)
	}

	res, err := client.DoWithAppSpan(req, "some_function")
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}
	res.Body.Close()

	// Output:
}

func someFunc(client *zipkinhttp.Client, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("some_function called with method: %s\n", r.Method)

		// retrieve span from context (created by server middleware)
		span := zipkin.SpanFromContext(r.Context())
		span.Tag("custom_key", "some value")

		// doing some expensive calculations....
		time.Sleep(25 * time.Millisecond)
		span.Annotate(time.Now(), "expensive_calc_done")

		newRequest, err := http.NewRequest("POST", url+"/other_function", nil)
		if err != nil {
			log.Printf("unable to create client: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}

		ctx := zipkin.NewContext(newRequest.Context(), span)

		newRequest = newRequest.WithContext(ctx)

		res, err := client.DoWithAppSpan(newRequest, "other_function")
		if err != nil {
			log.Printf("call to other_function returned error: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		res.Body.Close()
	}
}

func otherFunc(client *zipkinhttp.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("other_function called with method: %s\n", r.Method)
		time.Sleep(50 * time.Millisecond)
	}
}
