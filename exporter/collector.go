package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os/exec"
)

var _ prometheus.Collector = &Collector{}

type Collector struct {
	device       string
	PowerOnHours *prometheus.Desc
	Temperature  *prometheus.Desc
}

func NewCollector(device string) *Collector {
	var (
		labels = []string{
			"device",
			"model",
		}
	)
	return &Collector{
		device: device,
		PowerOnHours: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "power_on_hours"),
			"Power on hours",
			labels,
			nil,
		),
		Temperature: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "temperature"),
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
	if desc, err := c.collect(ch); err != nil {
		log.Printf("[ERROR] failed collecting metric %v: %v", desc, err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return
	}
}

func (c *Collector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	if c.device == "" {
		return nil, nil
	}

	out, err := exec.Command("smartctl", "-iA", c.device).CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] smart log: \n%s\n", out)
		return nil, err
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

	return nil, nil
}
