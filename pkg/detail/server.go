package detail

import (
	"log"
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/knative-sample/cloud-native-go-weather/pkg/db"
	"github.com/knative-sample/cloud-native-go-weather/pkg/tracing"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

type Server struct {
	Port             string
	ZipKinEndpoint   string
	ServiceName      string
	InstanceIp       string
	Beta             string
	TableStoreConfig *db.TableStoreConfig
	tracer           *zipkin.Tracer
}

func (wa *Server) Start() error {
	wa.tracer = tracing.GetTracer(wa.ServiceName, wa.InstanceIp, wa.ZipKinEndpoint)
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		wa.tracer, zipkinhttp.TagResponseSize(true),
	)

	router := mux.NewRouter()

	router.Use(serverMiddleware)
	//router.Use(utils.AccessLog)
	router.Methods("GET").Path("/api/area/weather/{adcode}/{date}").HandlerFunc(wa.GetDetail)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", wa.Port), router))

	return nil
}
