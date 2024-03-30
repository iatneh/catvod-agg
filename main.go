package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"tv-aggregation/app/conf"
	"tv-aggregation/app/models"
)

var (
	config *conf.Config
)

func init() {
	config = conf.AppConf
}
func main() {
	apiList, err := readJson()
	if err != nil {
		logrus.Errorf("read json file error: %s", err)
	}
	for i := range apiList {
		logrus.Infof("name:%s url:%s", apiList[i].Name, apiList[i].Url)
	}
}

func readJson() ([]models.ApiJson, error) {
	jsonContent, err := os.ReadFile(config.FileName)
	if err != nil {
		return nil, err
	}
	var apiList []models.ApiJson
	err = json.Unmarshal(jsonContent, &apiList)
	if err != nil {
		return nil, err
	}
	return apiList, nil
}
