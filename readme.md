# Exporter Backup Mysql

Exporter para backups de banco Mysql. Obtém o log do backup realizado e exporta como métricas para o Prometheus.

> Em desenvolvimento

## Estrutura de Log

A estrutura de log deve seguir o padrão:

```
timestamp | evento | origem | status | mensagem
```

Exemplo:

```
2024-09-08 20:30:01 | DUMP_INICIADO | database_1 | Inicio | Iniciando dump do banco de dados
2024-09-08 20:34:53 | DUMP_FINALIZADO | database_1 | OK | Dump do banco de dados finalizado
2024-09-08 20:34:53 | COMPRESSAO_INICIADA | database_1 | Inicio | Inicio compactação do arquivo
2024-09-08 20:56:25 | COMPRESSAO_CONCLUIDA | database_1 | OK | Arquivo compactado com sucesso
2024-09-08 20:56:25 | UPLOAD_INICIADO | database_1 | Inicio | Upload arquivo para o S3 iniciado
2024-09-08 20:56:46 | UPLOAD_CONCLUIDO | database_1 | OK | Upload realizado com sucesso
```

## Métricas

Exemplos de métricas que o exporter irá expor para o Prometheus:

```
dump_database_success_total{server="db-server",source="database_1"} 2
dump_database_success_total{server="db-server",source="database_2"} 2
dump_database_failed_total{server="db-server",source="database_1"} 1
dump_database_duration_seconds{reference="2024-09-08 20:30:01",server="db-server",source="database_1",start_time="2024-09-08 20:30:01"} 292
dump_database_duration_seconds{reference="2024-09-08 20:30:01",server="db-server",source="database_2",start_time="2024-09-08 20:56:46"} 10
dump_database_duration_seconds{reference="2024-09-08 21:30:01",server="db-server",source="database_1",start_time="2024-09-08 21:30:01"} 295
dump_database_duration_seconds{reference="2024-09-08 21:30:01",server="db-server",source="database_2",start_time="2024-09-08 21:56:57"} 11
```
