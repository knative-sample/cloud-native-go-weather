package weather

import (
	"github.com/knative-sample/cloud-native-go-weather/pkg/detail"
	"github.com/openzipkin/zipkin-go"
)

type WebApi struct {
	Port          string
	CityService   *Service
	DetailService *Service
	ResourceRoot  string

	ZipKinEndpoint string
	ServiceName    string
	InstanceIp     string
	tracer         *zipkin.Tracer
}

type Service struct {
	Host string
	Port string
}

type Detail struct {
	Detail []*detail.DetailInfo
}
