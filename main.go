package main

import (
	"flag"
	"github.com/io-developer/prom-smartctl-exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

var (
	listenAddr  = flag.String("listen", ":9167", "address for exporter")
	metricsPath = flag.String("path", "/metrics", "URL path for surfacing collected metrics")
)

func main() {
	flag.Parse()

	prometheus.MustRegister(exporter.New())

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Printf("starting exporter on %q", *listenAddr)

	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("cannot start exporter: %s", err)
	}
}
