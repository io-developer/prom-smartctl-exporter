package exporter

import (
	"log"
	"strings"
	"time"

	"github.com/io-developer/prom-smartctl-exporter/pkg/cmd"
	"github.com/prometheus/client_golang/prometheus"
)

const MIN_DEVICE_LEN = 5

type ExporterOpt struct {
	Shell          *cmd.Shell
	RescanInterval time.Duration
}

type Exporter struct {
	opt        ExporterOpt
	collectors []*Collector
}

func NewExporter(opt ExporterOpt) *Exporter {
	return &Exporter{
		opt:        opt,
		collectors: make([]*Collector, 0),
	}
}

func (e *Exporter) Start() (err error) {
	for {
		log.Printf("Device re-scanning..\n")
		err = e.rescan()
		if err != nil {
			return
		}
		log.Printf("Device collectors registered: %d\n", len(e.collectors))

		log.Printf("Device re-scanning wait for %.0f sec ...\n", e.opt.RescanInterval.Seconds())
		time.Sleep(e.opt.RescanInterval)
	}
}

func (e *Exporter) rescan() (err error) {
	devs, err := e.scanDevs()
	if err != nil {
		return
	}
	oldCollectors := e.collectors
	e.collectors = e.makeCollectors(devs)
	for _, col := range oldCollectors {
		prometheus.Unregister(col)
	}
	for _, col := range e.collectors {
		prometheus.MustRegister(col)
	}
	return
}

func (e *Exporter) scanDevs() (devs []string, err error) {
	stdout, _, _, err := e.opt.Shell.Exec("lsblk -Snpo name")
	if err != nil {
		return
	}
	rows := strings.Split(string(stdout), "\n")
	devs = make([]string, 0, len(rows))
	for _, row := range rows {
		if len(row) >= MIN_DEVICE_LEN {
			devs = append(devs, row)
		}
	}
	return
}

func (e *Exporter) makeCollectors(devs []string) []*Collector {
	cols := make([]*Collector, 0, len(devs))
	for _, dev := range devs {
		col := NewCollector(CollectorOpt{
			Device:        dev,
			Shell:         e.opt.Shell,
			SkipIfStandby: true,
		})
		if colErr := col.Validate(); colErr == nil {
			cols = append(cols, col)
		}
	}
	return cols
}
