package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aristidesneto/exporter-backup-mysql/metrics"
	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus"
)

func LoadFile(logPath string) {	
	log.Printf("Loading configuration file: %s\n", logPath)
	file, err := os.Open(logPath)
	if err != nil {
		log.Fatalf("Error to open file: %s", err)
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	for i, line := range lines {
		parserLogLine(i, line, lines)
    }
}

func parserLogLine(index int, line string, lines []string) {
	
	parts := strings.Split(line, "|")

	layoutDate := "2006-01-02 15:04:05"

	if len(parts) < 5 {
		log.Printf("This line is out of standard: %s", line)
		return
	}

	timestamp := strings.TrimSpace(parts[0])
	event := strings.TrimSpace(parts[1])
	source := strings.TrimSpace(parts[2])
	// status := strings.TrimSpace(parts[3])

	hostname := viper.GetString("server.hostname")

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
			metrics.M.DatabaseCounterSuccess.With(prometheus.Labels{"source": source, "server": hostname}).Inc()
			metrics.M.DatabaseDuration.WithLabelValues(source, hostname, backupStart, start_time.String()).Set(duration.Seconds())
		} else {
			metrics.M.DatabaseCounterFailed.With(prometheus.Labels{"source": source, "server": hostname}).Inc()
		}
	}
}