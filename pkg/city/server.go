package city

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/knative-sample/cloud-native-go-weather/pkg/db"
	"github.com/knative-sample/cloud-native-go-weather/pkg/tracing"
	zipkin "github.com/openzipkin/zipkin-go"
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
	router := mux.NewRouter()
	router.HandleFunc("/api/cities", wa.CityList)
	router.HandleFunc("/api/area/list/{citycode}", wa.AreaList)
	http.Handle("/", router)

	http.ListenAndServe(fmt.Sprintf(":%s", wa.Port), nil)

	return nil
}
