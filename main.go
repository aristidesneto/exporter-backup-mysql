package main

import (
	"bufio"
	"github.com/aristidesneto/exporter-backup-mysql/metrics"
	"github.com/aristidesneto/exporter-backup-mysql/parser"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)


func main()  {
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg)

	// Inicializar o Viper e configurar
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")

    err := viper.ReadInConfig()
    if err != nil {
        log.Printf("Erro ao ler o arquivo de configuração: %v", err)
    }

    logPath := viper.GetString("backup.log_path")
	hostname := viper.GetString("server.hostname")
	serverPort := viper.GetString("server.port")

	log.Printf("Arquivo de configuração carregado: %s\n", logPath)


	file, err := os.Open(logPath)
	if err != nil {
		log.Println("Erro ao abrir o arquivo:", err)
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
		parser.ParserLogLine(i, line, lines, *metrics, hostname)
    }
	
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Printf("Server is running on port %s", serverPort)
	http.ListenAndServe(":"+serverPort, nil)

}


