package city

import (
	"encoding/json"
	fmt "fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/openzipkin/zipkin-go"
)

type Citys struct {
	Citys []*City `protobuf:"bytes,1,rep,name=Citys,proto3" json:"Citys,omitempty"`
}

// The request message containing the user's name.
type City struct {
	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Citycode string `protobuf:"bytes,2,opt,name=citycode,proto3" json:"citycode,omitempty"`
}
type Areas struct {
	Areas []*Area `protobuf:"bytes,1,rep,name=Areas,proto3" json:"Areas,omitempty"`
}
type Area struct {
	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Citycode string `protobuf:"bytes,2,opt,name=citycode,proto3" json:"citycode,omitempty"`
}

func (s *Server) CityList(w http.ResponseWriter, r *http.Request) {
	if parent := zipkin.SpanFromContext(r.Context()); parent != nil {
		//tracer := tracing.GetTracer(s.ServiceName, s.InstanceIp, s.ZipKinEndpoint)
		subSpan := s.tracer.StartSpan("city_list", zipkin.Parent(parent.Context()))
		defer subSpan.Finish()
		//do some operations
		time.Sleep(time.Millisecond * 10)
	}
	cities, err := s.TableStoreConfig.QueryCities()
	if err != nil {
		log.Printf("QueryCities error %s", err.Error())
		fmt.Fprintf(w, err.Error())
		return
	}
	cs := make([]*City, 0)
	for _, city := range cities {
		c := &City{
			Name:     city.Name,
			Citycode: city.Citycode,
		}
		cs = append(cs, c)
	}
	dbts, _ := json.Marshal(&Citys{Citys: cs})
	fmt.Fprintf(w, string(dbts))
}

func (s *Server) AreaList(w http.ResponseWriter, r *http.Request) {
	if parent := zipkin.SpanFromContext(r.Context()); parent != nil {
		//tracer := tracing.GetTracer(s.ServiceName, s.InstanceIp, s.ZipKinEndpoint)
		subSpan := s.tracer.StartSpan("city_area_list", zipkin.Parent(parent.Context()))
		defer subSpan.Finish()
		//do some operations
		time.Sleep(time.Millisecond * 10)
	}
	params := strings.TrimPrefix(r.URL.Path[1:], "api/area/list/")
	vars := strings.Split(params, "/")
	citycode := vars[0]
	areas, err := s.TableStoreConfig.QueryAreaByCitycode(citycode)
	if err != nil {
		log.Printf("QueryAreaByCitycode error %s", err.Error())
		fmt.Fprintf(w, err.Error())
		return
	}
	as := make([]*Area, 0)
	for _, city := range areas {
		area := &Area{
			Name:     city.Name,
			Citycode: city.Adcode,
		}
		as = append(as, area)
	}
	dbts, _ := json.Marshal(&Areas{Areas: as})
	fmt.Fprintf(w, string(dbts))
}
