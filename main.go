package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/io-developer/prom-smartctl-exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddr    = flag.String("listen", ":9167", "address for exporter")
	metricsPath   = flag.String("path", "/metrics", "URL path for surfacing collected metrics")
	shellTemplate = flag.String("shell", "%s", "Shell template for system commands")
)

func main() {
	flag.Parse()

	shell := exporter.NewShell()
	shell.Template = *shellTemplate

	prometheus.MustRegister(exporter.New(shell))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Printf("starting exporter on %q", *listenAddr)

	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("cannot start exporter: %s", err)
	}
}
