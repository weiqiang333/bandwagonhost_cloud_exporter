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

var sum float64 = 1

func init() {
	pflag.String("exporter.address", ":9103", "The address on which to expose the web interface and generated Prometheus metrics.")
	pflag.String("config.file", "./config/bandwagonhost_cloud_exporter.yaml", "exporter config file")
}

type Exporter struct {
	MissingBlocks prometheus.Gauge
}

func NewExporter() *Exporter {
	return &Exporter{
		MissingBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "CapacityTotal",
			Help: "CapacityTotal",
		}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.MissingBlocks.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	sum = sum + 1

	infoMaps, err := pkg.GrabBwgServerInfo()
	if err != nil {
		return
	}
	fmt.Println(infoMaps)

	e.MissingBlocks.Set(sum)
	e.MissingBlocks.Collect(ch)
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
