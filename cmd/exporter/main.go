package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/io-developer/prom-smartctl-exporter/pkg/cmd"
	"github.com/io-developer/prom-smartctl-exporter/pkg/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddr    = flag.String("listen", ":9167", "address for exporter")
	shellTemplate = flag.String("shell", "%s", "Shell template for system commands")
)

func main() {
	flag.Parse()

	cmdShell := cmd.NewShell()
	cmdShell.Template = *shellTemplate

	exporter := exporter.NewExporter(cmdShell)
	err := exporter.Init()
	if err != nil {
		log.Fatalf("[ERROR] failed to init")
	}

	prometheus.MustRegister(exporter)

	log.Printf("starting exporter on %q", *listenAddr)

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("cannot start exporter: %s", err)
	}
}
