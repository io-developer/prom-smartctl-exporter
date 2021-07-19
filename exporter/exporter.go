package exporter

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	prometheus.Collector

	cmd        *Shell
	collectors []*Collector
}

func NewExporter(cmd *Shell) *Exporter {
	return &Exporter{
		cmd:        cmd,
		collectors: make([]*Collector, 0),
	}
}

func (e *Exporter) Init() error {
	e.collectors = make([]*Collector, 0)

	out, err := e.cmd.Exec("lsblk -Snpo name")
	if err == nil {
		devices := strings.Split(string(out), "\n")
		for _, device := range devices {
			collector := NewCollector(device, e.cmd)
			e.collectors = append(e.collectors, collector)
		}
	}
	return err
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range e.collectors {
		collector.Describe(ch)
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	for _, collector := range e.collectors {
		collector.Collect(ch)
	}
}
