package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/io-developer/prom-smartctl-exporter/pkg/cmd"
	"github.com/io-developer/prom-smartctl-exporter/pkg/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddr     = flag.String("listen", ":9167", "address for exporter")
	shellTemplate  = flag.String("shell", "%s", "Shell template for system commands")
	rescanInterval = flag.Duration("rescan", 300*time.Second, "Full device re-scan interval")
)

func main() {
	flag.Parse()

	cmdShell := cmd.NewShell()
	cmdShell.Template = *shellTemplate

	exporter := exporter.NewExporter(exporter.ExporterOpt{
		Shell:          cmdShell,
		RescanInterval: *rescanInterval,
	})

	go func() {
		err := exporter.Start()
		if err != nil {
			log.Fatalf("[ERROR] exporter error\n")
			panic(err)
		}
	}()

	log.Printf("starting exporter on %q", *listenAddr)

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("cannot start exporter: %s", err)
	}
}
