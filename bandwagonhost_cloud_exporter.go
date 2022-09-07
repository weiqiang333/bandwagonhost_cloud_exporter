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
	NodeSuspended prometheus.GaugeVec
}

func NewExporter() *Exporter {
	return &Exporter{
		NodeSuspended: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "node_suspended",
				Help: "server run status value, Running=0 / Suspended=1",
			},
			[]string{"ip_address_0", "node_ip", "hostname", "vm_type", "node_location", "os"}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.NodeSuspended.Describe(ch)
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
	}
	e.NodeSuspended.Collect(ch)
}

func main() {
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Println("Fatal error BindPFlags: %w", err.Error())
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(viper.GetString("config.file"))
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	prometheus.MustRegister(NewExporter())
	// http server
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8888", nil); err != nil {
		fmt.Println("Fatal error http: %w", err)
	}
}
