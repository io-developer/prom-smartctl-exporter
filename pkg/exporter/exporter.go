package exporter

import (
	"github.com/io-developer/prom-smartctl-exporter/pkg/cmd"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	prometheus.Collector

	shell      *cmd.Shell
	collectors []*Collector
}

func NewExporter(shell *cmd.Shell) *Exporter {
	return &Exporter{
		shell:      shell,
		collectors: make([]*Collector, 0),
	}
}

func (e *Exporter) Init() error {
	e.collectors = make([]*Collector, 0)

	// out, err := e.cmdShell.Exec("lsblk -Snpo name")
	// if err == nil {
	// 	devices := strings.Split(string(out), "\n")
	devices := []string{
		"/dev/sda",
	}
	for _, device := range devices {
		e.collectors = append(e.collectors, NewCollector(CollectorOpt{
			Device: device,
			Shell:  e.shell,
		}))
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
