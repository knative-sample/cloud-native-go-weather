package weather

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/knative-sample/cloud-native-go-weather/pkg/city"
	"github.com/openzipkin/zipkin-go"

	"encoding/json"

	"strings"

	"io/ioutil"

	"github.com/knative-sample/cloud-native-go-weather/pkg/detail"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

func (wa *WebApi) CityList(w http.ResponseWriter, r *http.Request) {

	currentSpan := wa.NewSpan("GetCityList", r.Context())
	defer currentSpan.Finish()

	addr := fmt.Sprintf("http://%s:%s/api/cities", wa.CityService.Host, wa.CityService.Port)
	if wa.tracer == nil {
		glog.Errorf("wa tracer is nill")
		fmt.Fprintf(w, "wa tracer is nill")
		return
	}
	client, err := zipkinhttp.NewClient(wa.tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		glog.Errorf("unable to create client: %+v\n", err)
	}
	res, _, err := SendReqest(client, "GetCityList", addr, currentSpan)
	if err != nil {
		glog.Errorf("CityList SendReqest error:%s", err.Error())
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, string(res))
}

func (wa *WebApi) Detail(w http.ResponseWriter, r *http.Request) {
	currentSpan := wa.NewSpan("GetDetail", r.Context())
	defer currentSpan.Finish()

	childSpan := wa.tracer.StartSpan("GetDetail", zipkin.Parent(currentSpan.Context()))
	defer childSpan.Finish()

	// 1. get city areas 2. foreach area get weather info
	params := strings.TrimPrefix(r.URL.Path[1:], "api/city/detail/")
	vars := strings.Split(params, "/")
	citycode := vars[0]
	date := vars[1]
	glog.Infof("citycode: %s", citycode)
	glog.Infof("date: %s", date)
	areaChildSpan := wa.tracer.StartSpan("GetDetail", zipkin.Parent(currentSpan.Context()))
	areas, err := wa.getAreas(citycode, areaChildSpan)
	if err != nil {
		glog.Errorf("getAreas error:%s", err.Error())
		return
	}
	defer areaChildSpan.Finish()

	detailResult := []*detail.DetailInfo{}
	for _, a := range areas {
		detailChildSpan := wa.tracer.StartSpan("GetDetail", zipkin.Parent(currentSpan.Context()))
		d, err := wa.getDetail(a.Citycode, date, detailChildSpan)
		if err != nil {
			glog.Errorf("getDetail error:%s", err.Error())
			continue
		}
		if d.Name == "" {
			continue
		}
		detailChildSpan.Finish()

		detailResult = append(detailResult, d)
	}

	dbts, _ := json.Marshal(detailResult)
	fmt.Fprintf(w, string(dbts))
}

func (wa *WebApi) getAreas(cityCode string, currentSpan zipkin.Span) ([]*city.Area, error) {
	//addr := "127.0.0.1:9090"
	addr := fmt.Sprintf("http://%s:%s/api/area/list/%s", wa.CityService.Host, wa.CityService.Port, cityCode)
	client, err := zipkinhttp.NewClient(wa.tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		glog.Errorf("unable to create client: %+v\n", err)
	}

	res, _, err := SendReqest(client, "getAreas", addr, currentSpan)
	if err != nil {
		glog.Errorf("getAreas SendReqest error:%s", err.Error())
		return nil, err
	}
	glog.Infof("getAreas areas: %s", res)
	areas := &city.Areas{}
	json.Unmarshal(res, areas)
	return areas.Areas, nil
}

func (wa *WebApi) getDetail(cityCode, date string, currentSpan zipkin.Span) (*detail.DetailInfo, error) {
	addr := fmt.Sprintf("http://%s:%s/api/area/weather/%s/%s", wa.DetailService.Host, wa.DetailService.Port, cityCode, date)
	client, err := zipkinhttp.NewClient(wa.tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		glog.Errorf("unable to create client: %+v\n", err)
	}

	res, _, err := SendReqest(client, "getDetail", addr, currentSpan)
	if err != nil {
		glog.Errorf("getDetail SendReqest error:%s", err.Error())
		return nil, err
	}
	detailInfo := &detail.DetailInfo{}
	json.Unmarshal(res, detailInfo)
	return detailInfo, nil
}

func SendReqest(client *zipkinhttp.Client, name, url string, currentSpan zipkin.Span) (body []byte, statusCode int, err error) {
	newRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Errorf("unable to create client: %+v\n", err)
		return
	}

	ctx := zipkin.NewContext(newRequest.Context(), currentSpan)

	newRequest = newRequest.WithContext(ctx)

	res, err := client.DoWithAppSpan(newRequest, name)
	if err != nil {
		glog.Errorf("call to other_function returned error: %+v\n", err)
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	statusCode = res.StatusCode
	//status code not in [200, 300) fail
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		fmt.Printf("response status code %d, error messge: %s", res.StatusCode, string(body))
		return
	}

	if err != nil {
		fmt.Printf("read the result of get url %s fails, response status code %d -- %v", url, res.StatusCode, err)
	}

	return
}