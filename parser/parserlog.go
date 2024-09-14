package parser

import (
	"github.com/aristidesneto/exporter-backup-mysql/metrics"
	"log"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func ParserLogLine(index int, line string, lines []string, metrics metrics.Metrics, hostname string) {
	
	parts := strings.Split(line, "|")

	layoutDate := "2006-01-02 15:04:05"

	if len(parts) < 4 {
		log.Println("Linha malformada:", line)
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
			log.Println(err)
		}
		
		end_time, err := time.Parse(layoutDate, strings.TrimSpace(strings.Split(lines[index + 1], "|")[0]))
		if err != nil {
			log.Println(err)
		}

		duration := end_time.Sub(start_time)
		
		// Pega o status da proxima linha
		nextLineEvent := strings.TrimSpace(strings.Split(nextLine, "|")[1])
		nextLineStatus := strings.TrimSpace(strings.Split(nextLine, "|")[3])

		if (nextLineEvent == "DUMP_FINALIZADO") && (nextLineStatus ==  "OK") {
			metrics.DatabaseCounterSuccess.With(prometheus.Labels{"source": source, "server": hostname}).Inc()
			metrics.DatabaseDuration.WithLabelValues(source, hostname, backupStart, start_time.String()).Set(duration.Seconds())
		} else {
			metrics.DatabaseCounterFailed.With(prometheus.Labels{"source": source, "server": hostname}).Inc()
		}
	}
}