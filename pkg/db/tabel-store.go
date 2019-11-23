package db

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search"
	"github.com/golang/glog"
)

const (
	isCityIndex   = "iscity"
	cityCodeIndex = "citycode"
)

type City struct {
	Citycode string `json:"citycode"`
	Name     string `json:"name"`
}
type Area struct {
	Adcode string `json:"adcode"`
	Name   string `json:"name"`
}
type Weather struct {
	Adcode       string `json:"adcode"`
	City         string `json:"city"`
	Date         string `json:"date"`
	Week         string `json:"week"`
	Dayweather   string `json:"dayweather"`
	Nightweather string `json:"nightweather"`
	Daytemp      string `json:"daytemp"`
	Nighttemp    string `json:"nighttemp"`
	Daywind      string `json:"daywind"`
	Nightwind    string `json:"nightwind"`
	Daypower     string `json:"daypower"`
	Nightpower   string `json:"nightpower"`
	Province     string `json:"province"`
	Reporttime   string `json:"reporttime"`
}

type TableStoreConfig struct {
	Endpoint        string
	TableName       string
	InstanceName    string
	AccessKeyId     string
	AccessKeySecret string
}

// 查询所有城市
func (cm *TableStoreConfig) QueryCities() ([]*City, error) {
	cities := make([]*City, 0)
	client := tablestore.NewClient(cm.Endpoint, cm.InstanceName, cm.AccessKeyId, cm.AccessKeySecret)
	searchRequest := &tablestore.SearchRequest{}
	searchRequest.SetTableName(cm.TableName)
	searchRequest.SetIndexName(isCityIndex)
	query := &search.MatchQuery{} // 设置查询类型为MatchQuery
	query.FieldName = isCityIndex // 设置要匹配的字段
	query.Text = "true"           // 设置要匹配的值
	searchQuery := search.NewSearchQuery()
	searchQuery.SetQuery(query)
	searchQuery.SetGetTotalCount(true)
	searchQuery.SetOffset(0)  //
	searchQuery.SetLimit(100) //
	searchRequest.SetSearchQuery(searchQuery)
	// 设置返回所有列
	searchRequest.SetColumnsToGet(&tablestore.ColumnsToGet{
		ReturnAll: true,
	})
	searchResponse, err := client.Search(searchRequest)
	if err != nil {
		fmt.Printf("%#v", err)
		return nil, err
	}
	for _, row := range searchResponse.Rows {
		cityMap := make(map[string]string, 0)
		if row.PrimaryKey.PrimaryKeys != nil {
			for _, col := range row.PrimaryKey.PrimaryKeys {
				cityMap[col.ColumnName] = col.Value.(string)
			}
		}
		for _, col := range row.Columns {
			cityMap[col.ColumnName] = col.Value.(string)
		}
		cb, err := json.Marshal(cityMap)
		if err != nil {
			fmt.Errorf("QueryCities Marshal error %s", err.Error())
			continue
		}
		city := &City{}
		err = json.Unmarshal(cb, city)
		if err != nil {
			fmt.Errorf("QueryCities Unmarshal error %s", err.Error())
			continue
		}
		cities = append(cities, city)
	}

	searchQuery = search.NewSearchQuery()
	searchQuery.SetQuery(query)
	searchQuery.SetGetTotalCount(true)
	searchQuery.SetOffset(100) //
	searchQuery.SetLimit(100)  //
	searchRequest.SetSearchQuery(searchQuery)
	searchRequest.SetColumnsToGet(&tablestore.ColumnsToGet{
		ReturnAll: true,
	})
	searchResponse, err = client.Search(searchRequest)
	if err != nil {
		fmt.Printf("%#v", err)
		return nil, err
	}
	for _, row := range searchResponse.Rows {
		cityMap := make(map[string]string, 0)
		if row.PrimaryKey.PrimaryKeys != nil {
			for _, col := range row.PrimaryKey.PrimaryKeys {
				cityMap[col.ColumnName] = col.Value.(string)
			}
		}
		for _, col := range row.Columns {
			cityMap[col.ColumnName] = col.Value.(string)
		}
		cb, err := json.Marshal(cityMap)
		if err != nil {
			fmt.Errorf("QueryCities Marshal error %s", err.Error())
			continue
		}
		city := &City{}
		err = json.Unmarshal(cb, city)
		if err != nil {
			fmt.Errorf("QueryCities Unmarshal error %s", err.Error())
			continue
		}
		cities = append(cities, city)
	}

	return cities, nil
}

// 根据城市代码查询区域
func (cm *TableStoreConfig) QueryAreaByCitycode(citycode string) ([]*Area, error) {
	areas := make([]*Area, 0)
	client := tablestore.NewClient(cm.Endpoint, cm.InstanceName, cm.AccessKeyId, cm.AccessKeySecret)
	searchRequest := &tablestore.SearchRequest{}
	searchRequest.SetTableName(cm.TableName)
	searchRequest.SetIndexName(cityCodeIndex)
	query := &search.MatchQuery{}   // 设置查询类型为MatchQuery
	query.FieldName = cityCodeIndex // 设置要匹配的字段
	query.Text = citycode           // 设置要匹配的值
	searchQuery := search.NewSearchQuery()
	searchQuery.SetQuery(query)
	searchQuery.SetGetTotalCount(true)
	searchQuery.SetOffset(0)  //
	searchQuery.SetLimit(100) //
	searchRequest.SetSearchQuery(searchQuery)
	// 设置返回所有列
	searchRequest.SetColumnsToGet(&tablestore.ColumnsToGet{
		ReturnAll: true,
	})
	searchResponse, err := client.Search(searchRequest)
	if err != nil {
		fmt.Printf("%#v", err)
		return nil, err
	}
	for _, row := range searchResponse.Rows {
		iscity := false
		areaMap := make(map[string]string, 0)
		if row.PrimaryKey.PrimaryKeys != nil {
			for _, col := range row.PrimaryKey.PrimaryKeys {
				areaMap[col.ColumnName] = col.Value.(string)
			}
		}
		for _, col := range row.Columns {
			if col.ColumnName == "iscity" && col.Value.(string) == "true" {
				iscity = true
			}
			areaMap[col.ColumnName] = col.Value.(string)
		}
		if iscity {
			continue
		}
		cb, err := json.Marshal(areaMap)
		if err != nil {
			fmt.Errorf("QueryAreaByCitycode Marshal error %s", err.Error())
			continue
		}
		area := &Area{}
		err = json.Unmarshal(cb, area)
		if err != nil {
			fmt.Errorf("QueryAreaByCitycode Unmarshal error %s", err.Error())
			continue
		}
		areas = append(areas, area)
	}
	return areas, nil
}

// 查询区域天气
// - adcode: 区域代码
// - date: 查询日期。格式：2019-09-26
//
func (cm *TableStoreConfig) QueryWeather(adcode, date string) (*Weather, error) {
	client := tablestore.NewClient(cm.Endpoint, cm.InstanceName, cm.AccessKeyId, cm.AccessKeySecret)
	getRowRequest := new(tablestore.GetRowRequest)
	criteria := new(tablestore.SingleRowQueryCriteria)
	putPk := &tablestore.PrimaryKey{}

	putPk.AddPrimaryKeyColumn("adcode", adcode)
	putPk.AddPrimaryKeyColumn("date", date)
	criteria.PrimaryKey = putPk

	getRowRequest.SingleRowQueryCriteria = criteria
	getRowRequest.SingleRowQueryCriteria.TableName = cm.TableName
	getRowRequest.SingleRowQueryCriteria.MaxVersion = 1
	getResp, err := client.GetRow(getRowRequest)
	if err != nil {
		glog.Errorf("QueryWeather failed with error: %s", err.Error())
		return nil, err
	}
	weatherMap := make(map[string]string, 0)
	if getResp.PrimaryKey.PrimaryKeys != nil {
		for _, col := range getResp.PrimaryKey.PrimaryKeys {
			weatherMap[col.ColumnName] = col.Value.(string)
		}
	}
	if getResp.Columns != nil {
		for _, col := range getResp.Columns {
			// 过滤掉 id 信息
			if col.ColumnName == "id" {
				continue
			}
			weatherMap[col.ColumnName] = col.Value.(string)
		}
	}
	cb, err := json.Marshal(weatherMap)
	if err != nil {
		fmt.Errorf("QueryWeather Marshal error %s", err.Error())
		return nil, err
	}
	weather := &Weather{}
	err = json.Unmarshal(cb, weather)
	if err != nil {
		fmt.Errorf("QueryWeather Unmarshal error %s", err.Error())
		return nil, err
	}
	return weather, nil
}
