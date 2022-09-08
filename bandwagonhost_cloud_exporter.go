package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"bandwagonhost_cloud_exporter/pkg"
)

func init() {
	pflag.String("exporter.address", ":9103", "The address on which to expose the web interface and generated Prometheus metrics.")
	pflag.String("config.file", "./config/bandwagonhost_cloud_exporter.yaml", "exporter config file")
}

type Exporter struct {
	NodeSuspended         prometheus.GaugeVec
	PlanMonthlyData       prometheus.GaugeVec
	PlanMonthlyDataGb     prometheus.GaugeVec
	MonthlyDataMultiplier prometheus.GaugeVec
	DataCounter           prometheus.GaugeVec
	DataCounterGb         prometheus.GaugeVec
	AvailableTrafficGb    prometheus.GaugeVec
	DataNextReset         prometheus.GaugeVec
}

func NewExporter() *Exporter {
	return &Exporter{
		NodeSuspended: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "node_suspended",
				Help: "server run status value, Running=0 / Suspended=1",
			}, []string{"ip_address", "node_ip", "hostname", "vm_type", "node_location", "os"}),
		PlanMonthlyData: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "plan_monthly_data",
				Help: "每月可用流量 (bytes)",
			}, []string{"ip_address", "hostname"}),
		PlanMonthlyDataGb: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "plan_monthly_data_gb",
				Help: "每月可用流量 (GB)",
			}, []string{"ip_address", "hostname"}),
		MonthlyDataMultiplier: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "monthly_data_multiplier",
				Help: "宽带流量计费系数",
			}, []string{"ip_address", "hostname"}),
		DataCounter: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "data_counter",
				Help: "当月已用流量 (bytes)",
			}, []string{"ip_address", "hostname"}),
		DataCounterGb: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "data_counter_gb",
				Help: "当月已用流量 (GB)",
			}, []string{"ip_address", "hostname"}),
		AvailableTrafficGb: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "available_traffic_gb",
				Help: "剩余可用流量 (GB)",
			}, []string{"ip_address", "hostname"}),
		DataNextReset: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "data_next_reset",
				Help: "流量计数器重置的日期和时间（UNIX 时间戳）",
			}, []string{"ip_address", "hostname"}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.NodeSuspended.Describe(ch)
	e.PlanMonthlyData.Describe(ch)
	e.PlanMonthlyDataGb.Describe(ch)
	e.MonthlyDataMultiplier.Describe(ch)
	e.DataCounter.Describe(ch)
	e.DataCounterGb.Describe(ch)
	e.AvailableTrafficGb.Describe(ch)
	e.DataNextReset.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	infoMaps, err := pkg.GrabBwgServerInfo()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, infoMap := range infoMaps {
		var nodeSuspendedCode float64 = 0
		if infoMap.Suspended {
			nodeSuspendedCode = 1
		}
		e.NodeSuspended.With(prometheus.Labels{
			"ip_address":    infoMap.IpAddresses[0],
			"node_ip":       infoMap.Nodeip,
			"hostname":      infoMap.Hostname,
			"vm_type":       infoMap.VmType,
			"node_location": infoMap.NodeLocation,
			"os":            infoMap.OS,
		}).Set(nodeSuspendedCode)
		e.PlanMonthlyData.WithLabelValues(infoMap.IpAddresses[0], infoMap.Hostname).Set(infoMap.PlanMonthlyData)
		e.PlanMonthlyDataGb.WithLabelValues(infoMap.IpAddresses[0], infoMap.Hostname).Set(infoMap.PlanMonthlyDataGb)
		e.MonthlyDataMultiplier.WithLabelValues(infoMap.IpAddresses[0], infoMap.Hostname).Set(infoMap.MonthlyDataMultiplier)
		e.DataCounter.WithLabelValues(infoMap.IpAddresses[0], infoMap.Hostname).Set(infoMap.DataCounter)
		e.DataCounterGb.WithLabelValues(infoMap.IpAddresses[0], infoMap.Hostname).Set(infoMap.DataCounterGb)
		e.AvailableTrafficGb.WithLabelValues(infoMap.IpAddresses[0], infoMap.Hostname).Set(infoMap.AvailableTrafficGb)
		e.DataNextReset.WithLabelValues(infoMap.IpAddresses[0], infoMap.Hostname).Set(infoMap.DataNextReset)
	}
	e.NodeSuspended.Collect(ch)
	e.PlanMonthlyData.Collect(ch)
	e.PlanMonthlyDataGb.Collect(ch)
	e.MonthlyDataMultiplier.Collect(ch)
	e.DataCounter.Collect(ch)
	e.DataCounterGb.Collect(ch)
	e.AvailableTrafficGb.Collect(ch)
	e.DataNextReset.Collect(ch)
}

func main() {
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Println("Fatal error BindPFlags: %w", err.Error())
	}
	fmt.Println("load config file ", viper.GetString("config.file"))
	viper.SetConfigType("yaml")
	viper.SetConfigFile(viper.GetString("config.file"))
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	prometheus.MustRegister(NewExporter())
	// http server
	fmt.Printf("http server start, address %s/metrics\n", viper.GetString("exporter.address"))
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(viper.GetString("exporter.address"), nil); err != nil {
		fmt.Println("Fatal error http: %w", err)
	}
}
