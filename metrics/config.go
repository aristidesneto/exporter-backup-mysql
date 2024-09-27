package metrics

import (
	"log"
	"os"
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
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

func PushMetrics(metric string)  {
	url_pushgateway := os.Getenv("PUSHGATEWAY_URL")
	auth_username := os.Getenv("PUSHGATEWAY_AUTH_USER")
	auth_password := os.Getenv("PUSHGATEWAY_AUTH_PASS")

	if url_pushgateway == "" {
		log.Fatalln("Missing PUSHGATEWAY_URL variable")
	}

	if auth_username == "" {
		log.Fatalln("Missing PUSHGATEWAY_AUTH_USER variable")
	}

	if auth_password == "" {
		log.Fatalln("Missing PUSHGATEWAY_AUTH_PASS variable")
	}
	pusher := push.New(url_pushgateway, "mysql_backup").BasicAuth(auth_username, auth_password)

	// Acessando o campo dinamicamente
    field, err := getMetricByName(M, metric)
    if err != nil {
        log.Println(err)
        return
    }

	if err := pusher.Collector(field).
		Grouping("instance", metric).
		Push(); err != nil {
		log.Fatalf("Error to send metrics: %v", err)
	}

	log.Printf("Metrics %s send successfully", metric)	
}

// Função para acessar dinamicamente o campo da struct usando uma string
func getMetricByName(m *Metrics, fieldName string) (prometheus.Collector, error) {
    v := reflect.ValueOf(m).Elem() // Obtém o valor da struct referenciada

    field := v.FieldByName(fieldName) // Busca o campo pelo nome da string
    if !field.IsValid() {
        log.Printf("Field %s not found", fieldName)
    }

	collector, ok := field.Interface().(prometheus.Collector)
	if !ok {
		log.Println("Field not implement collector", field)
	}

    return collector, nil
}