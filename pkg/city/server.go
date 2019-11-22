package city

import (
	fmt "fmt"
	"log"
	"net/http"

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
	TableStoreConfig *db.TableStoreConfig
	tracer           *zipkin.Tracer
}

func (wa *Server) Start() error {
	wa.tracer = tracing.GetTracer(wa.ServiceName, wa.InstanceIp, wa.ZipKinEndpoint)
	// create global zipkin http server middleware
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		wa.tracer, zipkinhttp.TagResponseSize(true),
	)

	// initialize router
	router := mux.NewRouter()
	router.Use(serverMiddleware)

	router.Methods("GET").Path("/api/cities").HandlerFunc(wa.CityList)
	router.Methods("GET").Path("/api/area/list/{citycode}").HandlerFunc(wa.AreaList)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", wa.Port), router))

	return nil
}
