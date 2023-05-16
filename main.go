package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/songjiayang/nginx-log-exporter/collector"
	"github.com/songjiayang/nginx-log-exporter/config"
)

var (
	bind, configFile string
)

func main() {
	flag.StringVar(&bind, "web.listen-address", ":9147", "Address to listen on for the web interface and API.")
	flag.StringVar(&configFile, "config.file", "config.yml", "Nginx log exporter configuration file name.")
	flag.Parse()

	cfg, err := config.LoadFile(configFile)
	if err != nil {
		log.Panic(err)
	}

	for _, app := range cfg.App {
		go collector.NewCollector(app).Run()
	}

	fmt.Printf("running HTTP server on address %s\n", bind)

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Fatalf("start server with error: %v\n", err)
	}
}
