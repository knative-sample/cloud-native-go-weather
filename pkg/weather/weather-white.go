package weather

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/knative-sample/cloud-native-go-weather/pkg/city"

	"encoding/json"

	"strings"

	"io/ioutil"

	"github.com/knative-sample/cloud-native-go-weather/pkg/detail"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

func (wa *WebApi) DetailWhite(w http.ResponseWriter, r *http.Request) {

	// 1. get city areas 2. foreach area get weather info
	params := strings.TrimPrefix(r.URL.Path[1:], "api/city/detail/")
	vars := strings.Split(params, "/")
	citycode := vars[0]
	date := vars[1]
	glog.Infof("citycode: %s", citycode)
	glog.Infof("date: %s", date)
	areas, err := wa.getAreasWhite(citycode)
	if err != nil {
		glog.Errorf("getAreas error:%s", err.Error())
		return
	}

	detailResult := []*detail.DetailInfo{}
	for _, a := range areas {
		d, err := wa.getDetailWhite(a.Citycode, date)
		if err != nil {
			glog.Errorf("getDetail error:%s", err.Error())
			continue
		}
		if d.Name == "" {
			continue
		}

		detailResult = append(detailResult, d)
	}

	dbts, _ := json.Marshal(detailResult)
	fmt.Fprintf(w, string(dbts))
}

func (wa *WebApi) getAreasWhite(cityCode string) ([]*city.Area, error) {
	//addr := "127.0.0.1:9090"
	addr := fmt.Sprintf("http://%s:%s/api/area/list/%s", wa.CityService.Host, wa.CityService.Port, cityCode)
	client, err := zipkinhttp.NewClient(wa.tracer, zipkinhttp.ClientTrace(false))
	if err != nil {
		glog.Errorf("unable to create client: %+v\n", err)
	}

	res, _, err := SendReqestWhite(client, "getAreas", addr)
	if err != nil {
		glog.Errorf("getAreas SendReqest error:%s", err.Error())
		return nil, err
	}
	glog.Infof("getAreas areas: %s", res)
	areas := &city.Areas{}
	json.Unmarshal(res, areas)
	return areas.Areas, nil
}

func (wa *WebApi) getDetailWhite(cityCode, date string) (*detail.DetailInfo, error) {
	addr := fmt.Sprintf("http://%s:%s/api/area/weather/%s/%s", wa.DetailService.Host, wa.DetailService.Port, cityCode, date)
	client, err := zipkinhttp.NewClient(wa.tracer, zipkinhttp.ClientTrace(false))
	if err != nil {
		glog.Errorf("unable to create client: %+v\n", err)
	}

	res, _, err := SendReqestWhite(client, "getDetail", addr)
	if err != nil {
		glog.Errorf("getDetail SendReqest error:%s", err.Error())
		return nil, err
	}
	detailInfo := &detail.DetailInfo{}
	json.Unmarshal(res, detailInfo)
	return detailInfo, nil
}

func SendReqestWhite(client *zipkinhttp.Client, name, url string) (body []byte, statusCode int, err error) {
	newRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Errorf("unable to create client: %+v\n", err)
		return
	}

	res, err := client.Do(newRequest)
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
