# Exporter Mysql Backup

Exporter to Mysql database backups. Gets the log of the backup performed and exports it as metrics to Prometheus. 

> Working In Progress

## Log structure

The log structure must follow the pattern:

```
timestamp | event | source | status | message
```

Example:

```
2024-09-08 20:30:01 | DUMP_INICIADO | database_1 | Inicio | Iniciando dump do banco de dados
2024-09-08 20:34:53 | DUMP_FINALIZADO | database_1 | OK | Dump do banco de dados finalizado
2024-09-08 20:34:53 | COMPRESSAO_INICIADA | database_1 | Inicio | Inicio compactação do arquivo
2024-09-08 20:56:25 | COMPRESSAO_CONCLUIDA | database_1 | OK | Arquivo compactado com sucesso
2024-09-08 20:56:25 | UPLOAD_INICIADO | database_1 | Inicio | Upload arquivo para o S3 iniciado
2024-09-08 20:56:46 | UPLOAD_CONCLUIDO | database_1 | OK | Upload realizado com sucesso
```

## Metrics

Examples of metrics that the exporter will exponse to Prometheus:

```
dump_database_success_total{server="db-server",source="database_1"} 2
dump_database_success_total{server="db-server",source="database_2"} 2
dump_database_failed_total{server="db-server",source="database_1"} 1
dump_database_duration_seconds{reference="2024-09-08 20:30:01",server="db-server",source="database_1",start_time="2024-09-08 20:30:01"} 292
dump_database_duration_seconds{reference="2024-09-08 20:30:01",server="db-server",source="database_2",start_time="2024-09-08 20:56:46"} 10
dump_database_duration_seconds{reference="2024-09-08 21:30:01",server="db-server",source="database_1",start_time="2024-09-08 21:30:01"} 295
dump_database_duration_seconds{reference="2024-09-08 21:30:01",server="db-server",source="database_2",start_time="2024-09-08 21:56:57"} 11
```


## Systemd config

```conf
[Unit]
Description=Exporter Mysql Backup Log
After=network.target
 
[Service]
Type=simple
User=root
Group=root
ExecStart=/etc/exporter-mysql-backup/exporter --logpath=/var/log/backups/database_dump.log
WorkingDirectory=/etc/exporter-mysql-backup
 
[Install]
WantedBy=multi-user.target
```

## Crontab

```
# Exporter mysql backup
20 * * * * user /path/to/script/start.sh >> /path/to/script/exporter.log 2>&1
```