package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	databaseCounterSuccess *prometheus.CounterVec
	databaseCounterFailed *prometheus.CounterVec
	databaseDuration *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *metrics  {
	m := &metrics{
		databaseCounterSuccess: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "dump_database_success_total",
				Help: "Total backup database successfully",
			},
			[]string{"source", "server"},
		),
		databaseCounterFailed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "dump_database_failed_total",
				Help: "Total backup database failed",
			},
			[]string{"source", "server"},
		),
		databaseDuration: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "dump_database_duration_seconds",
				Help: "Duration of the backup",
			},
			[]string{"source", "server", "reference", "start_time"},
		),
	}

	reg.MustRegister(m.databaseCounterSuccess)
	reg.MustRegister(m.databaseCounterFailed)
	reg.MustRegister(m.databaseDuration)
	return m
}

func main()  {
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	logFile := "backup.log"

	file, err := os.Open(logFile)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo: ", err)
		return
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	for i, line := range lines {
		processLogLine(i, line, lines, *m)
    }
	
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	fmt.Println("Server is up: port 8888")
	http.ListenAndServe(":8888", nil)

}

func processLogLine(index int, line string, lines []string, m metrics) {
	
	parts := strings.Split(line, "|")

	layoutDate := "2006-01-02 15:04:05"

	if len(parts) < 4 {
		fmt.Println("Linha malformada:", line)
		return
	}

	timestamp := strings.TrimSpace(parts[0])
	event := strings.TrimSpace(parts[1])
	source := strings.TrimSpace(parts[2])
	// status := strings.TrimSpace(parts[3])

	if event == "DUMP_INICIADO" {
		nextLine := lines[index + 1]
		backupStart := strings.TrimSpace(parts[0])


		start_time, err := time.Parse(layoutDate, timestamp)
		if err != nil {
			fmt.Println(err)
		}
		
		end_time, err := time.Parse(layoutDate, strings.TrimSpace(strings.Split(lines[index + 1], "|")[0]))
		if err != nil {
			fmt.Println(err)
		}

		duration := end_time.Sub(start_time)
		
		// Pega o status da proxima linha
		nextLineEvent := strings.TrimSpace(strings.Split(nextLine, "|")[1])
		nextLineStatus := strings.TrimSpace(strings.Split(nextLine, "|")[3])
		serverName := "srv-aiguilles-banco"

		if (nextLineEvent == "DUMP_FINALIZADO") && (nextLineStatus ==  "OK") {
			m.databaseCounterSuccess.With(prometheus.Labels{"source": source, "server": serverName}).Inc()
			m.databaseDuration.WithLabelValues(source, serverName, backupStart, start_time.String()).Set(duration.Seconds())
		} else {
			m.databaseCounterFailed.With(prometheus.Labels{"source": source, "server": serverName}).Inc()
		}
	}
}
