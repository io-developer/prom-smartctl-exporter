package exporter

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var BaseLabels = []string{
	"device",
	"model",
}

var AttrLabels = append(
	BaseLabels,
	"id",
	"name",
	"flag",
	"flag_string",
	"flag_prefailure",
	"flag_updated_online",
	"flag_performance",
	"flag_error_rate",
	"flag_event_count",
	"flag_auto_keep",
)

type Collector struct {
	device       string
	cmd          *Shell
	PowerOnHours *prometheus.Desc
	Temperature  *prometheus.Desc
	AttrValue    *prometheus.Desc
	AttrWorst    *prometheus.Desc
	AttrThresh   *prometheus.Desc
	AttrRaw      *prometheus.Desc
}

func NewCollector(device string, cmd *Shell) *Collector {
	return &Collector{
		device: device,
		cmd:    cmd,
		PowerOnHours: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "smartctl_power_on_hours"),
			"Power on hours",
			BaseLabels,
			nil,
		),
		Temperature: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "smartctl_temperature"),
			"Temperature",
			BaseLabels,
			nil,
		),
		AttrValue: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "smartctl_attr_value"),
			"S.M.A.R.T. attribute value",
			AttrLabels,
			nil,
		),
		AttrWorst: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "smartctl_attr_worst"),
			"S.M.A.R.T. attribute worst",
			AttrLabels,
			nil,
		),
		AttrThresh: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "smartctl_attr_thresh"),
			"S.M.A.R.T. attribute threshold",
			AttrLabels,
			nil,
		),
		AttrRaw: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "smartctl_attr_raw"),
			"S.M.A.R.T. attribute raw value",
			AttrLabels,
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

	baseLabels := []string{
		c.device,
		smart.GetInfo("Device Model", "Model Family"),
	}

	ch <- prometheus.MustNewConstMetric(
		c.PowerOnHours,
		prometheus.GaugeValue,
		float64(smart.GetAttr(9).rawValue),
		baseLabels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.Temperature,
		prometheus.GaugeValue,
		float64(smart.GetAttr(190, 194).rawValue),
		baseLabels...,
	)

	attrLabels := append(
		baseLabels,
		"177",
		"Wear_Leveling_Count",
		"19",
		"PO--C- ",
		"1",
		"0",
		"0",
		"0",
		"1",
		"0",
	)
	ch <- prometheus.MustNewConstMetric(
		c.AttrValue,
		prometheus.GaugeValue,
		float64(98),
		attrLabels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.AttrWorst,
		prometheus.GaugeValue,
		float64(98),
		attrLabels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.AttrThresh,
		prometheus.GaugeValue,
		float64(0),
		attrLabels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.AttrRaw,
		prometheus.GaugeValue,
		float64(43),
		attrLabels...,
	)
}
