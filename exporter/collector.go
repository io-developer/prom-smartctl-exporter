package exporter

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &Collector{}

type Collector struct {
	device       string
	cmd          *Shell
	PowerOnHours *prometheus.Desc
	Temperature  *prometheus.Desc
}

func NewCollector(device string, cmd *Shell) *Collector {
	var (
		labels = []string{
			"device",
			"model",
		}
	)
	return &Collector{
		device: device,
		cmd:    cmd,
		PowerOnHours: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "smartctl_power_on_hours"),
			"Power on hours",
			labels,
			nil,
		),
		Temperature: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "smartctl_temperature"),
			"Temperature",
			labels,
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.PowerOnHours,
		c.Temperature,
	}
	for _, d := range ds {
		ch <- d
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	if c.device == "" {
		return
	}

	out, err := c.cmd.Exec(fmt.Sprintf("smartctl -iA %s", c.device))
	if err != nil {
		log.Printf("[ERROR] smart log: \n%s\n", out)
		return
	}

	smart := ParseSmart(string(out))
	labels := []string{
		c.device,
		smart.GetInfo("Device Model", "Model Family"),
	}

	ch <- prometheus.MustNewConstMetric(
		c.PowerOnHours,
		prometheus.GaugeValue,
		float64(smart.GetAttr(9).rawValue),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.Temperature,
		prometheus.GaugeValue,
		float64(smart.GetAttr(190, 194).rawValue),
		labels...,
	)
}
