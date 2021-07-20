package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	prometheus.Collector

	cmdShell   *CmdShell
	collectors []*Collector
}

func NewExporter(cmdShell *CmdShell) *Exporter {
	return &Exporter{
		cmdShell:   cmdShell,
		collectors: make([]*Collector, 0),
	}
}

func (e *Exporter) Init() error {
	e.collectors = make([]*Collector, 0)

	// out, err := e.cmdShell.Exec("lsblk -Snpo name")
	// if err == nil {
	// 	devices := strings.Split(string(out), "\n")
	for _, device := range []string{"/dev/sda"} {
		collector := NewCollector(device, e.cmdShell)
		e.collectors = append(e.collectors, collector)
	}
	// }
	return nil //err
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
