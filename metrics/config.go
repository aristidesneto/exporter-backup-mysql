package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)


type Metrics struct {
	DatabaseCounterSuccess *prometheus.CounterVec
	DatabaseCounterFailed *prometheus.CounterVec
	DatabaseDuration *prometheus.GaugeVec
}

var M *Metrics

func NewMetrics(reg *prometheus.Registry) *Metrics  {
	metrics := &Metrics{
		DatabaseCounterSuccess: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "dump_database_success_total",
				Help: "Total backup database successfully",
			},
			[]string{"source", "server"},
		),
		DatabaseCounterFailed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "dump_database_failed_total",
				Help: "Total backup database failed",
			},
			[]string{"source", "server"},
		),
		DatabaseDuration: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "dump_database_duration_seconds",
				Help: "Duration of the backup",
			},
			[]string{"source", "server", "reference", "start_time"},
		),
	}

	// Verificar se as métricas foram criadas corretamente
	if metrics.DatabaseCounterSuccess == nil || metrics.DatabaseCounterFailed == nil || metrics.DatabaseDuration == nil {
		log.Fatal("Failed to initialize one or more metrics")
	}

	reg.MustRegister(metrics.DatabaseCounterSuccess)
	reg.MustRegister(metrics.DatabaseCounterFailed)
	reg.MustRegister(metrics.DatabaseDuration)

	M = metrics

	return metrics
}