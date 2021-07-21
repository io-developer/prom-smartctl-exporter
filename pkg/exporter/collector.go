package exporter

import (
	"fmt"
	"log"

	"github.com/io-developer/prom-smartctl-exporter/pkg/cmd"
	"github.com/io-developer/prom-smartctl-exporter/pkg/data"
	"github.com/prometheus/client_golang/prometheus"
)

type CollectorOpt struct {
	Device string
	Shell  *cmd.Shell
}

func (o CollectorOpt) GetConstLabels() prometheus.Labels {
	out, err := o.Shell.Exec(fmt.Sprintf("smartctl -iA -l scttempsts --json=ou %s", o.Device))
	if err != nil {
		log.Printf("[ERROR] smart log: \n%s\n", out)
		return prometheus.Labels{}
	}

	resp, err := data.ParseSmartctlJson(out)
	if err != nil {
		log.Printf("[ERROR] parse smartctl json: \n%v\n", err)
		return prometheus.Labels{}
	}

	log.Print(resp)

	return prometheus.Labels{
		"device": "SomeDevice",
		"model":  "SomeModel",
	}
}

type Collector struct {
	opt          CollectorOpt
	PowerOnHours prometheus.Gauge
	Temperature  prometheus.Gauge
	AttrValue    *prometheus.GaugeVec
	AttrWorst    *prometheus.GaugeVec
	AttrThresh   *prometheus.GaugeVec
	AttrRaw      *prometheus.GaugeVec
}

func NewCollector(opt CollectorOpt) *Collector {
	constLabels := opt.GetConstLabels()
	attrLabelNames := []string{
		"id",
		"name",
		"is_prefailure",
		"is_updated_online",
		"is_performance",
		"is_error_rate",
		"is_event_count",
		"is_auto_keep",
	}
	return &Collector{
		opt: opt,

		PowerOnHours: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "smartctl_power_on_hours",
			Help:        "Power on hours",
			ConstLabels: constLabels,
		}),
		Temperature: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "smartctl_temperature",
			Help:        "Temperature",
			ConstLabels: constLabels,
		}),
		AttrValue: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   "",
				Subsystem:   "",
				Name:        "smartctl_attr_value",
				Help:        "S.M.A.R.T. attribute value",
				ConstLabels: constLabels,
			},
			attrLabelNames,
		),
		AttrWorst: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   "",
				Subsystem:   "",
				Name:        "smartctl_attr_worst",
				Help:        "S.M.A.R.T. attribute worst",
				ConstLabels: constLabels,
			},
			attrLabelNames,
		),
		AttrThresh: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   "",
				Subsystem:   "",
				Name:        "smartctl_attr_thresh",
				Help:        "S.M.A.R.T. attribute threshold",
				ConstLabels: constLabels,
			},
			attrLabelNames,
		),
		AttrRaw: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   "",
				Subsystem:   "",
				Name:        "smartctl_attr_raw",
				Help:        "S.M.A.R.T. attribute raw value",
				ConstLabels: constLabels,
			},
			attrLabelNames,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.PowerOnHours.Describe(ch)
	c.Temperature.Describe(ch)
	c.AttrValue.Describe(ch)
	c.AttrWorst.Describe(ch)
	c.AttrThresh.Describe(ch)
	c.AttrRaw.Describe(ch)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	if c.opt.Device == "" {
		return
	}

	out, err := c.opt.Shell.Exec(fmt.Sprintf("smartctl -n standby -iA %s", c.opt.Device))
	if err != nil {
		log.Printf("[ERROR] smart log: \n%s\n", out)
		return
	}

	smart := OldParseSmart(string(out))

	c.PowerOnHours.Set(float64(smart.GetAttr(9).rawValue))
	c.PowerOnHours.Collect(ch)

	c.Temperature.Set(float64(smart.GetAttr(190, 194).rawValue))
	c.Temperature.Collect(ch)

	attrLabels := prometheus.Labels{
		"id":                "177",
		"name":              "Wear_Leveling_Count",
		"is_prefailure":     "true",
		"is_updated_online": "false",
		"is_performance":    "false",
		"is_error_rate":     "false",
		"is_event_count":    "true",
		"is_auto_keep":      "false",
	}

	c.AttrValue.With(attrLabels).Set(float64(98))
	c.AttrValue.Collect(ch)

	c.AttrWorst.With(attrLabels).Set(float64(98))
	c.AttrWorst.Collect(ch)

	c.AttrThresh.With(attrLabels).Set(float64(0))
	c.AttrThresh.Collect(ch)

	c.AttrRaw.With(attrLabels).Set(float64(43))
	c.AttrRaw.Collect(ch)
}
