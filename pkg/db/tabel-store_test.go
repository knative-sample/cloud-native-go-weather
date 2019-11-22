package db

import (
	"encoding/json"
	"testing"
)

func TestQueryCities(t *testing.T) {
	ts := &TableStoreConfig{
		Endpoint:        "https://knative-xxx.cn-beijing.ots.aliyuncs.com",
		TableName:       "city",
		InstanceName:    "knative-weather",
		AccessKeyId:     "xxxx",
		AccessKeySecret: "xxxx",
	}
	cities, err := ts.QueryCities()
	if err != nil {
		t.Errorf("Test QueryCities error %s", err.Error())
		return
	}
	cb, _ := json.Marshal(cities)
	t.Logf("cities: %s", cb)
}

func TestQueryAreaByCitycode(t *testing.T) {
	ts := &TableStoreConfig{
		Endpoint:        "https://knative-xxx.cn-beijing.ots.aliyuncs.com",
		TableName:       "city",
		InstanceName:    "knative-weather",
		AccessKeyId:     "xxxx",
		AccessKeySecret: "xxxx",
	}
	areas, err := ts.QueryAreaByCitycode("010")
	if err != nil {
		t.Errorf("Test QueryAreaByCitycode error %s", err.Error())
		return
	}
	cb, _ := json.Marshal(areas)
	t.Logf("areas: %s", cb)
}

func TestQueryWeather(t *testing.T) {
	ts := &TableStoreConfig{
		Endpoint:        "https://knative-xxx.cn-beijing.ots.aliyuncs.com",
		TableName:       "weather",
		InstanceName:    "knative-weather",
		AccessKeyId:     "xxxx",
		AccessKeySecret: "xxxx",
	}
	weather, err := ts.QueryWeather("110101", "2019-11-11")
	if err != nil {
		t.Errorf("Test QueryWeather error %s", err.Error())
		return
	}
	cb, _ := json.Marshal(weather)
	t.Logf("weather: %s", cb)
}
