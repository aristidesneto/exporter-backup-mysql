package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/aristidesneto/exporter-backup-mysql/config"
	"github.com/aristidesneto/exporter-backup-mysql/metrics"
	"github.com/aristidesneto/exporter-backup-mysql/parser"
	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var reg *prometheus.Registry

func init()  {
	// Configuration file
	config.Configuration("./config")

	// Prometheus metrics
	reg = prometheus.NewRegistry()
	log.Println("Registry initialized")
	metrics.NewMetrics(reg)
	if metrics.M == nil {
		log.Fatal("Metrics not initialized")
	} else {
		log.Println("Metrics initialized successfully")
	}
}


func main()  {
	var logPath string
	flag.StringVar(&logPath, "logpath", "backup.log", "Informe o caminho do arquivo de backup de log")
	flag.Parse()

	if reg == nil {
		log.Fatal("Registry is nil before setting up HTTP handler")
	}

	// Loading backup file
	parser.LoadFile(logPath)

	serverPort := viper.GetString("server.port")
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Printf("Server is running on port %s", serverPort)
	err := http.ListenAndServe(":" + serverPort, nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}