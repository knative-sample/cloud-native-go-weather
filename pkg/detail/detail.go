package detail

import (
	"encoding/json"
	fmt "fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/knative-sample/cloud-native-go-weather/pkg/utils/logs"
	"github.com/openzipkin/zipkin-go"
)

type DetailInfo struct {
	Adcode       string `protobuf:"bytes,1,opt,name=adcode,proto3" json:"adcode,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Date         string `protobuf:"bytes,3,opt,name=date,proto3" json:"date,omitempty"`
	Daypower     string `protobuf:"bytes,4,opt,name=daypower,proto3" json:"daypower,omitempty"`
	Daytemp      string `protobuf:"bytes,5,opt,name=daytemp,proto3" json:"daytemp,omitempty"`
	Dayweather   string `protobuf:"bytes,6,opt,name=dayweather,proto3" json:"dayweather,omitempty"`
	Daywind      string `protobuf:"bytes,7,opt,name=daywind,proto3" json:"daywind,omitempty"`
	Nightpower   string `protobuf:"bytes,8,opt,name=nightpower,proto3" json:"nightpower,omitempty"`
	Nighttemp    string `protobuf:"bytes,9,opt,name=nighttemp,proto3" json:"nighttemp,omitempty"`
	Nightweather string `protobuf:"bytes,10,opt,name=nightweather,proto3" json:"nightweather,omitempty"`
	Nightwind    string `protobuf:"bytes,11,opt,name=nightwind,proto3" json:"nightwind,omitempty"`
	Province     string `protobuf:"bytes,12,opt,name=province,proto3" json:"province,omitempty"`
	Reporttime   string `protobuf:"bytes,13,opt,name=reporttime,proto3" json:"reporttime,omitempty"`
	Week         string `protobuf:"bytes,14,opt,name=week,proto3" json:"week,omitempty"`
	Limit        string `protobuf:"bytes,15,opt,name=limit,proto3" json:"limit"`
}

// SayHello implements helloworld.GreeterServer
func (s *Server) GetDetail(w http.ResponseWriter, r *http.Request) {
	params := strings.TrimPrefix(r.URL.Path[1:], "api/area/weather/")
	vars := strings.Split(params, "/")
	citycode := vars[0]
	date := vars[1]
	if parent := zipkin.SpanFromContext(r.Context()); parent != nil {
		//tracer := tracing.GetTracer(s.ServiceName, s.InstanceIp, s.ZipKinEndpoint)
		subSpan := s.tracer.StartSpan("detail_sub_span", zipkin.Parent(parent.Context()))
		defer subSpan.Finish()
		//do some operations
		time.Sleep(time.Millisecond * 10)
	}
	wi, err := s.TableStoreConfig.QueryWeather(citycode, date)
	if err != nil {
		log.Printf("QueryWeather error %s", err.Error())
		fmt.Fprintf(w, err.Error())
		return
	}
	d := &DetailInfo{
		Adcode:       citycode,
		Name:         wi.City,
		Date:         wi.Date,
		Daypower:     wi.Daypower,
		Daytemp:      wi.Daytemp,
		Dayweather:   wi.Dayweather,
		Daywind:      wi.Daywind,
		Nightpower:   wi.Nightpower,
		Nighttemp:    wi.Nighttemp,
		Nightweather: wi.Nightweather,
		Nightwind:    wi.Nightwind,
		Province:     wi.Province,
		Reporttime:   wi.Reporttime,
		Week:         wi.Week,
		Limit:        "",
	}
	if s.Beta == "true" {
		d.Limit = fmt.Sprintf("今日限行尾号：5，0")
	}
	dbts, _ := json.Marshal(d)
	l := &logs.Log{}
	l.Info("DETAIL", "Get weather detail", r)
	fmt.Fprintf(w, string(dbts))

}
