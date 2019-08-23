Prometheus
==========

Prometheus é um sistema open-source de monitoramento e alerta. Este sistema coleta métricas e processa regras, armazenando-os em um banco de dados de série temporal (_time series database_ **tsdb**).


Tipos de Métrica
----------------

Os tipos de métrica foram explicados na sessão anterior **[1\_metrics]**.


Configurações
-------------

Para o Prometheus fazer o _scrape_ de um alvo, é preciso configurá-lo um pouco<sup>[configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)</sup>. É preciso dizer o endereço do alvo, se o esquema é http ou https, de quanto em quanto tempo é para ser feito a coleta, qual o tempo que ele pode gastar processando as regras, quanto tempo pode ficar fazendo o _scrape_ de um único alvo, e assim por diante. Essa configuração é definida pelo arquivo **prometheus.yml**. Os parâmetros são passados através do **docker-compose.yml** e o **startup.sh** constroi o arquivo de configuração. Eis um exemplo do arquivo:
**prometheus.yml**
```
global:
  scrape_interval: 30s
  evaluation_interval: 15s
  scrape_timeout: 10s

rule_files: 
  - aggr_global_http_per_server.yml

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
    - targets: ['localhost:9090']

  - job_name: "metrics_generator/metrics"
    metrics_path: /metrics
    file_sd_configs:
      - files:
        - targets.json
```

A seção `global` define as configurações que serão usadas para todos os `jobs`. O `scrape_interval` define de quanto em quanto tempo o Prometheus vai coletar as métricas nos alvos (este intervalo de tempo independe do tempo que levou para fazer a coleta<sup>[prometheus-ticker](https://utcc.utoronto.ca/~cks/space/blog/sysadmin/PrometheusScrapeIntervalBit)</sup>). O `evaluation_interval` é o tempo que o Prometheus tem para processar as regras definidas nos arquivos de regras. O `scrape_timeout` limita o tempo que o Prometheus pode fazer a coleta, por alvo. A seção seguinte é a `rule_files` no qual possui uma lista de arquivos de regras. Cada arquivo `yml` embaixo possui o nome do grupo das regras, o nome da regra e a expressão para processar antes de armazenar no tsdb.

Na seção `scrape_configs` estão definidas os alvos para coletar as métricas. Para cada alvo, é preciso um nome, um caminho de métricas (_endpoint_), e um arquivo com a lista dos alvos (pode ser um `dns:porta` ou `ip:porta`). No caso deste exemplo, como está sendo usado o `docker-compose` para executar os serviços, o alvo no `targets.json` pode ser com o nome do serviço (`gerador:3000`) pois o docker resolve o nome e encontra seu endereço.


rules
-----

As regras de gravação são definidas de acordo com a estrutura definida no site 
<sup>[rules](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)</sup>

```
groups:
- name: http_requests_duration_seconds
  rules:
  - record: http_requests_duration_seconds_sum
    expr: sum(irate(http_requests_duration_seconds_sum[1m])) by (status, uri)

  - record: http_requests_duration_seconds_count
    expr: sum(irate(http_requests_duration_seconds_count[1m])) by (status, uri)

  - record: http_requests_duration_seconds_average
    expr:
      sum(irate(http_requests_duration_seconds_sum[1m])) by (status, uri)
      /
      sum(irate(http_requests_duration_seconds_count[1m])) by (status, uri)
```


alerts
------

<sup>[alerts](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)</sup>




**docker-compose.yml**
```
version: '3.3'

services:

  prometheus:
    image: mfurquim/prometheus:1.0.0
    build:
      context: .
      dockerfile: ./Dockerfile
      args:
        scrapeinterval: 30s
        evaluationinterval: 15s
        scrapetimeout: 10s
        tsdbretention: 3d
        targetsfile: targets.json
        rulespath: rules
        alertspath: alerts
        scheme: http
        metricspaths: /metrics,/metrics-http,/metrics-negocio
        startupfile: startup.sh
    ports:
      - 9090:9090
    volumes:
      - prometheus:/prometheus

  generator:
    image: mfurquim/metrics-generator:v1.0.0
    build: ./metrics_generator/
    ports:
      - 3000:3000

volumes:
  prometheus:
```

**Dockerfile**
```
FROM prom/prometheus:v2.12.0

#### ARGS #####

# Defines the path o the files for startup script, targets, rules, and alerts
# - startupfile is the script to initialize the container
# - targetsfile contains a list of targets to scrape
# - rulespath is a directory which contains all the files for record rules
# - alertspath is a directory which contains all the files for alert rules
# - metricspaths is the endpoint which the metrics are exposed to scrape
ARG startupfile
ARG targetsfile
ARG rulespath
ARG alertspath
ARG metricspaths

# Defines the configuration of the Prometheus instance
# - scrapeinterval is the interval in seconds that it will collect the metrics
# - evaluationinterval is the time in seconds that it has to process the record rules and store them
# - scrapetimeout is the time in seconds for each scrape to timeout
# - tsdbretention is how long it should keep data in the database
# - scheme is either http or https
ARG scrapeinterval
ARG evaluationinterval
ARG scrapetimeout
ARG tsdbretention
ARG scheme

#### ENVS ####

ENV SCRAPE_INTERVAL ${scrapeinterval}
ENV EVALUATION_INTERVAL ${evaluationinterval}
ENV SCRAPE_TIMEOUT ${scrapetimeout}
ENV TSDB_RETENTION ${tsdbretention}
ENV SCHEME ${scheme}
ENV METRICS_PATHS ${metricspaths}

#### CONFIG ####

USER root

ADD $targetsfile /etc/prometheus/targets.json
ADD $rulespath /etc/prometheus/
ADD $alertspath /etc/prometheus/
ADD $startupfile /

ADD prometheus.yml /etc/prometheus/
ADD build.sh /

RUN chmod -R 755 /etc/prometheus/
RUN chmod -R 755 /startup.sh
RUN chmod +x /build.sh

RUN sh /build.sh /etc/prometheus/

ENTRYPOINT [ "/bin/sh" ]
CMD [ "/startup.sh" ]
```
