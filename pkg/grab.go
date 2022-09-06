package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func getUrl() (urls []string, err error) {
	getinfoUrl := viper.GetString("bandwagonhost.getinfo_url")
	serverApiKeys, ok := viper.Get("bandwagonhost.server_api_key").([]interface{})
	if ok == false {
		return urls, fmt.Errorf("faile in read bandwagonhost.server_api_key type error of config file")
	}
	for i := 0; i < len(serverApiKeys); i++ {
		serverApiKey := viper.GetStringMap(fmt.Sprintf("bandwagonhost.server_api_key.%v", i))
		veid, key := serverApiKey["veid"], serverApiKey["key"]
		urls = append(urls, fmt.Sprintf("%s?veid=%v&api_key=%s", getinfoUrl, veid, key))
	}
	return urls, nil
}

func createData(infoMap map[string]interface{}) map[string]interface{} {
	monthly_data_multiplier := infoMap["monthly_data_multiplier"].(float64)
	plan_monthly_data := infoMap["plan_monthly_data"].(float64) * monthly_data_multiplier
	data_counter := infoMap["data_counter"].(float64) * monthly_data_multiplier
	plan_monthly_data_gb := formatSizeConversion(plan_monthly_data)
	data_counter_gb := formatSizeConversion(data_counter)
	available_traffic_gb := formatSizeConversion(plan_monthly_data - data_counter)
	plan_ram_gb := formatSizeConversion(infoMap["plan_ram"].(float64))
	plan_disk_gb := formatSizeConversion(infoMap["plan_disk"].(float64))
	return map[string]interface{}{
		"vm_type":              infoMap["vm_type"].(string),
		"hostname":             infoMap["hostname"].(string),
		"node_ip":              infoMap["node_ip"].(string),
		"node_location":        infoMap["node_location"].(string),
		"plan_monthly_data_gb": plan_monthly_data_gb,
		"data_counter_gb":      data_counter_gb,
		"available_traffic_gb": available_traffic_gb,
		"plan_ram_gb":          plan_ram_gb,
		"plan_disk_gb":         plan_disk_gb,
		"os":                   infoMap["os"].(string),
		"data_next_reset":      infoMap["data_next_reset"].(float64),
		"suspended":            infoMap["suspended"].(bool),
	}
}

// formatSizeConversion: bytes -> GB
func formatSizeConversion(size float64) float64 {
	sizeGB, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", size/float64(1024)/float64(1024)/float64(1024)), 64)
	return sizeGB
}

// GrabBwgServerInfo: grab bwg cloud server info
// plan_monthly_data: 每月可用流量 (bytes), 基于 monthly_data_multiplier 系数.
// data_counter: 当月已用流量 (bytes), 基于 monthly_data_multiplier 系数.
// monthly_data_multiplier: 宽带流量计费系数, 与此相乘.
func GrabBwgServerInfo() ([]map[string]interface{}, error) {
	var infoMaps []map[string]interface{}
	urls, err := getUrl()
	if err != nil {
		fmt.Println("failed in GrabBwgServerInfo:", err.Error())
		return infoMaps, fmt.Errorf("faile in GrabBwgServerInfo: %w", err)
	}
	for _, url := range urls {
		client := http.Client{
			Timeout: 2 * time.Second,
		}
		resp, err := client.Get(url)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("failed in GrabBwgServerInfo: %w", err.Error())
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		var infoMap map[string]interface{}
		if err := json.Unmarshal(body, &infoMap); err != nil {
			fmt.Println("failed in GrabBwgServerInfo json Unmarshal: %w", err.Error())
			continue
		}
		newInfoMap := createData(infoMap)
		infoMaps = append(infoMaps, newInfoMap)
	}
	return infoMaps, nil
}
