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

type ServerInfo struct {
	VmType                string   `json:"vm_type"`
	Hostname              string   `json:"hostname"`
	Nodeip                string   `json:"node_ip"`
	IpAddresses           []string `json:"ip_addresses"`
	NodeLocation          string   `json:"node_location"`
	MonthlyDataMultiplier float64  `json:"monthly_data_multiplier"`
	PlanMonthlyData       float64  `json:"plan_monthly_data"`
	DataCounter           float64  `json:"data_counter"`
	PlanMonthlyDataGb     float64  `json:"plan_monthly_data_gb"`
	DataCounterGb         float64  `json:"data_counter_gb"`
	AvailableTrafficGb    float64  `json:"available_traffic_gb"`
	PlanRam               float64  `json:"plan_ram"`
	PlanDisk              float64  `json:"plan_disk"`
	PlanRamGb             float64  `json:"plan_ram_gb"`
	PlanDiskGb            float64  `json:"plan_disk_gb"`
	OS                    string   `json:"os"`
	DataNextReset         float64  `json:"data_next_reset"`
	Suspended             bool     `json:"suspended"`
	Error                 float64  `json:"error"`
	Message               string   `json:"message"`
}

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

func createData(infoMap ServerInfo) ServerInfo {
	plan_monthly_data := infoMap.PlanMonthlyData * infoMap.MonthlyDataMultiplier
	data_counter := infoMap.DataCounter * infoMap.MonthlyDataMultiplier
	plan_monthly_data_gb := formatSizeConversion(plan_monthly_data)
	data_counter_gb := formatSizeConversion(data_counter)
	available_traffic_gb := formatSizeConversion(plan_monthly_data - data_counter)
	plan_ram_gb := formatSizeConversion(infoMap.PlanRam)
	plan_disk_gb := formatSizeConversion(infoMap.PlanDisk)

	infoMap.PlanMonthlyDataGb = plan_monthly_data_gb
	infoMap.DataCounterGb = data_counter_gb
	infoMap.AvailableTrafficGb = available_traffic_gb
	infoMap.PlanRamGb = plan_ram_gb
	infoMap.PlanDiskGb = plan_disk_gb
	return infoMap
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
func GrabBwgServerInfo() ([]ServerInfo, error) {
	var infoMaps []ServerInfo
	urls, err := getUrl()
	if err != nil {
		fmt.Println("failed in GrabBwgServerInfo: ", err.Error())
		return infoMaps, fmt.Errorf("faile in GrabBwgServerInfo: %w", err)
	}
	for _, url := range urls {
		client := http.Client{
			Timeout: 3 * time.Second,
		}
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("failed in GrabBwgServerInfo: ", err.Error())
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		var infoMap ServerInfo
		if err := json.Unmarshal(body, &infoMap); err != nil {
			fmt.Println("failed in GrabBwgServerInfo json Unmarshal: ", err.Error(), string(body))
			continue
		}
		if infoMap.Error != 0 {
			fmt.Printf("failed in GrabBwgServerInfo body error code %v error message %v, url %s\n", infoMap.Error, infoMap.Message, url)
			continue
		}
		newInfoMap := createData(infoMap)
		infoMaps = append(infoMaps, newInfoMap)
	}
	return infoMaps, nil
}
