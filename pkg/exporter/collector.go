package exporter

import (
	"fmt"
	"log"
	"strconv"

	"github.com/io-developer/prom-smartctl-exporter/pkg/cmd"
	"github.com/io-developer/prom-smartctl-exporter/pkg/data"
	"github.com/prometheus/client_golang/prometheus"
)

type CollectorOpt struct {
	Device        string
	Shell         *cmd.Shell
	SkipIfStandby bool
}

func (o CollectorOpt) getConstLabels() (labels prometheus.Labels) {
	resp, _, err := o.getSmartctlResponse("-i -l scttempsts")
	if err != nil {
		return
	}
	return prometheus.Labels{
		"model":    resp.ModelName,
		"serial":   resp.SerialNumber,
		"firmware": resp.FirmwareVersion,
	}
}

func (o CollectorOpt) getSmartctlResponse(cmdOpts string) (resp data.Response, exitCode int, err error) {
	stdout, stderr, exitCode, err := o.Shell.Exec(
		fmt.Sprintf("smartctl %s --json=ou %s", cmdOpts, o.Device),
	)
	if err != nil && exitCode != 0 && exitCode != 2 {
		log.Printf("[ERROR] smartctl: %#v\n%#v\n", err, stderr)
		return
	}
	resp, err = data.ParseSmartctlJson(stdout)
	if err != nil {
		log.Printf("[ERROR] smartctl parse: %#v\n", err)
	}
	return
}

type Collector struct {
	opt             CollectorOpt
	hasFirstCollect bool
	PowerOnHours    prometheus.Gauge
	Temperature     prometheus.Gauge
	AttrValue       *prometheus.GaugeVec
	AttrWorst       *prometheus.GaugeVec
	AttrThresh      *prometheus.GaugeVec
	AttrRaw         *prometheus.GaugeVec
}

func NewCollector(opt CollectorOpt) *Collector {
	constLabels := opt.getConstLabels()
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

	cmdOpts := "-n standby,7 -iA -l scttempsts"
	if c.hasFirstCollect && c.opt.SkipIfStandby {
		cmdOpts = "-n standby,7 -iA -l scttempsts"
	}
	c.hasFirstCollect = true

	resp, exitCode, err := c.opt.getSmartctlResponse(cmdOpts)
	if exitCode == 7 {
		resp, _, err = c.opt.getSmartctlResponse("-i -l scttempsts")
	}
	if err != nil {
		return
	}

	c.updateMetrics(resp)

	c.PowerOnHours.Collect(ch)
	c.Temperature.Collect(ch)
	c.AttrValue.Collect(ch)
	c.AttrWorst.Collect(ch)
	c.AttrThresh.Collect(ch)
	c.AttrRaw.Collect(ch)
}

func (c *Collector) updateMetrics(resp data.Response) {
	if resp.PowerOnTime.Hours != 0 {
		c.PowerOnHours.Set(float64(resp.PowerOnTime.Hours))
	}
	if resp.Temperature.Current != 0 {
		c.Temperature.Set(float64(resp.Temperature.Current))
	}
	if resp.AtaSctStatus.DeviceState.String == "Active" {
		for _, attr := range resp.AtaSmartSttributes.Table {
			attrLabels := prometheus.Labels{
				"id":                strconv.FormatInt(int64(attr.Id), 10),
				"name":              attr.Name,
				"is_prefailure":     strconv.FormatBool(attr.Flags.Prefailure),
				"is_updated_online": strconv.FormatBool(attr.Flags.UpdatedOnline),
				"is_performance":    strconv.FormatBool(attr.Flags.Performance),
				"is_error_rate":     strconv.FormatBool(attr.Flags.ErrorRate),
				"is_event_count":    strconv.FormatBool(attr.Flags.EventCount),
				"is_auto_keep":      strconv.FormatBool(attr.Flags.AutoKeep),
			}
			c.AttrValue.With(attrLabels).Set(float64(attr.Value))
			c.AttrWorst.With(attrLabels).Set(float64(attr.Worst))
			c.AttrThresh.With(attrLabels).Set(float64(attr.Thresh))
			c.AttrRaw.With(attrLabels).Set(float64(attr.Raw.Value))
		}
	}
}
