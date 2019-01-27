package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os/exec"
	"strings"
)

const (
	namespace = "smartctl"
)

// An Exporter is a Prometheus exporter for metrics.
// It wraps all metrics collectors and provides a single global
// exporter which can serve metrics.
//
// It implements the exporter.Collector interface in order to register
// with Prometheus.
type Exporter struct {
}

var _ prometheus.Collector = &Exporter{}

// New creates a new Exporter which collects metrics by creating a apcupsd
// client using the input ClientFunc.
func New() *Exporter {
	return &Exporter{}
}

// Describe sends all the descriptors of the collectors included to
// the provided channel.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	NewCollector("").Describe(ch)
}

// Collect sends the collected metrics from each of the collectors to
// exporter.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	out, err := exec.Command("lsblk", "-Snpo", "name").CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] failed collecting metric %v: %v", out, err)
		return
	}
	devices := strings.Split(string(out), "\n")
	for _, device := range devices {
		NewCollector(device).Collect(ch)
	}
}
