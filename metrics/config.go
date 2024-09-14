package metrics

import "github.com/prometheus/client_golang/prometheus"


type Metrics struct {
	DatabaseCounterSuccess *prometheus.CounterVec
	DatabaseCounterFailed *prometheus.CounterVec
	DatabaseDuration *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics  {
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

	reg.MustRegister(metrics.DatabaseCounterSuccess)
	reg.MustRegister(metrics.DatabaseCounterFailed)
	reg.MustRegister(metrics.DatabaseDuration)

	return metrics
}