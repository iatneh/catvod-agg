package aggregation

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/marcozac/go-jsonc"
	"github.com/sirupsen/logrus"
	"os"
	"tv-aggregation/app/conf"
	"tv-aggregation/app/models"
)

func AggToFile() {
	multiApiConfig, err := aggregationMultiSiteInfo()
	if err != nil {
		logrus.Errorf("agg site info err:%s", err)
		return
	}
	for _, singleApiConfig := range multiApiConfig.Urls {
		logrus.Debugf("api name:%s url:%s", singleApiConfig.Name, singleApiConfig.Url)
	}
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(multiApiConfig)

	err = os.WriteFile(conf.AppConf.General.ToFilePath, bf.Bytes(), 0666)
	if err != nil {
		logrus.Errorf("write json file err:%s", err)
	}
}
func readJsonFile() ([]models.SingleApiConfig, error) {
	jsonContent, err := os.ReadFile(conf.AppConf.General.FileName)
	if err != nil {
		return nil, err
	}
	var apiList []models.SingleApiConfig
	err = json.Unmarshal(jsonContent, &apiList)
	if err != nil {
		return nil, err
	}
	return apiList, nil
}

func aggregationMultiSiteInfo() (*models.MultiApiConfig, error) {
	sites, err := readJsonFile()
	if err != nil {
		return nil, err
	}
	var totalList = models.MultiApiConfig{
		Urls: make([]models.SingleApiConfig, 0),
	}
	for _, site := range sites {
		err := getSiteApiConfig(&totalList, site)
		if err != nil {
			continue
		}
	}
	return &totalList, nil
}

func getSiteApiConfig(totalList *models.MultiApiConfig, site models.SingleApiConfig) error {
	if site.Type == "single" {
		err := urlCheck(site.Url)
		if err != nil {
			logrus.Warnf("check url fail,site:%s url name:%s url:%s %s",
				site.Name, site.Name, site.Url, err)
			return err
		}
		totalList.Urls = append(totalList.Urls, site)
		return nil
	}
	var apiContent models.MultiApiConfig
	client := resty.New()
	resp, err := client.R().Get(site.Url)
	if err != nil {
		logrus.Errorf("call api %s error,%s", site.Url, err)
		return err
	}
	respBody, err := jsonc.Sanitize(resp.Body())
	if err != nil {
		logrus.Errorf("discard comment in %s resp error,%s", site.Url, err)
	}
	err = json.Unmarshal(bytes.TrimPrefix(respBody, []byte("\xef\xbb\xbf")), &apiContent)
	if err != nil {
		logrus.Errorf("unmarshal json in %s resp error,resp data:%s,%s", site.Url, respBody, err)
		return err
	}
	for _, urlDetail := range apiContent.Urls {
		if searchUrlExists(totalList, urlDetail) {
			continue
		}
		// 测试链接是否有效
		err := urlCheck(urlDetail.Url)
		if err != nil {
			logrus.Warnf("check url fail,site:%s url name:%s url:%s %s",
				site.Name, urlDetail.Name, urlDetail.Url, err)
			continue
		}
		urlDetail.Name = site.Name + "-" + urlDetail.Name
		totalList.Urls = append(totalList.Urls, urlDetail)
	}
	return nil
}

// urlCheck 测试链接是否有效
func urlCheck(url string) error {
	client := resty.New()
	_, err := client.R().Get(url)
	return err
}

func searchUrlExists(multiApiConfig *models.MultiApiConfig, singleApiConfig models.SingleApiConfig) bool {
	urls := multiApiConfig.Urls
	for _, url := range urls {
		if url.Url == singleApiConfig.Url {
			return true
		}
	}
	return false
}
