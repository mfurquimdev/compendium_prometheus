#!/bin/bash

#WARNING: This code has ugly implementations because Prometheus base image doesn't have bash, just sh!

echo "Generating prometheus.yml according to ENV variables..."
FILE=/etc/prometheus/prometheus.yml

#global
cat > $FILE <<- EOM
global:
  scrape_interval: $SCRAPE_INTERVAL
  evaluation_interval: $EVALUATION_INTERVAL
  scrape_timeout: $SCRAPE_TIMEOUT

EOM

RULES=""
NEWLINE=$'\n'
for file in /etc/prometheus/*.yml; do
    FILENAME="$(expr $file \: '/etc/prometheus/\(.*\)')"
    if [ ! $FILENAME == "prometheus.yml" ]; then
        RULES="${RULES}${NEWLINE}  - ${FILENAME}"
    fi
done

cat >> $FILE <<- EOM
rule_files: $RULES

EOM


#alert managers
if [ "$ALERTMANAGER_TARGETS" != "" ]; then
    cat >> $FILE <<- EOM
alerting:
  alertmanagers:
  - static_configs:
    - targets:
EOM
    #add each alert manager target
    for i in $(echo $ALERTMANAGER_TARGETS | tr " " "\n")
    do
    cat >> $FILE <<- EOM
      - $i
EOM
    done
fi


cat >> $FILE <<- EOM
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
    - targets: ['localhost:9090']

EOM

#static scrapers
if [ "$STATIC_SCRAPE_TARGETS" != "" ]; then
    #add each static scrape target
    for SL in $(echo $STATIC_SCRAPE_TARGETS | tr " " "\n")
    do
        #this has to be done this ugly way because we don't have bash here, just sh!
        NAME=''
        HOST=''
        METRICS_PATH=''
        i=0
        for ST in $(echo $SL | tr "@" "\n")
        do
          if [ $i -eq 0 ]; then
            NAME=$ST
            i=1
          else
            HOST=$ST
          fi
        done
        
        METRICS_PATH=$(echo $HOST | cut -d/ -f2-)
        echo $METRICS_PATH
        if [ "$METRICS_PATH" == "" ] || [ "$METRICS_PATH" == "$HOST" ]; then
          METRICS_PATH="metrics"
        fi
        HOST=$(echo $HOST | cut -d/ -f1)

        echo $SCHEME_SCRAPE_TARGETS
         if [ "$SCHEME_SCRAPE_TARGETS" == "http" ] || [ "$SCHEME_SCRAPE_TARGETS" == "" ] ; then
          SCHEME_SCRAPE_TARGETS="http"
        fi

         if [ "$SCHEME_SCRAPE_TARGETS" == "https" ]; then
          SCHEME_SCRAPE_TARGETS="https"
          TLS_IGNORE="tls_config:"
          TRUE="insecure_skip_verify: true"

        fi

        cat >> $FILE <<- EOM      
 
        
  - job_name: '$NAME'
    metrics_path: /$METRICS_PATH
    scheme: $SCHEME_SCRAPE_TARGETS
    $TLS_IGNORE
       $TRUE
    static_configs:
    - targets: ['$HOST']

EOM
    done
fi

#dns scrapers
if [ "$DNS_SCRAPE_TARGETS" != "" ]; then
    #add each static scrape target
    for SL in $(echo $DNS_SCRAPE_TARGETS | tr " " "\n")
    do
        #this has to be done this ugly way because we don't have bash here, just sh!
        NAME=''
        HOSTPORT=''
        PORT=''
        HOST=''
        METRICS_PATH=''
        a=0
        for ST in $(echo $SL | tr "@" "\n")
        do
          if [ $a -eq 0 ]; then
            NAME=$ST
            a=1
          else
            HOSTPORT=$ST
          fi
        done

        METRICS_PATH=$(echo $HOSTPORT | cut -d/ -f2-)
        echo $METRICS_PATH
        if [ "$METRICS_PATH" == "" ] || [ "$METRICS_PATH" == "$HOSTPORT" ]; then
          METRICS_PATH="metrics"
        fi
        HOSTPORT=$(echo $HOSTPORT | cut -d/ -f1)

        for HP in $(echo $HOSTPORT | tr ":" "\n")
        do
          if [ $a -eq 1 ]; then
            HOST=$HP
            a=2
          else
            PORT=$HP
          fi
        done
        cat >> $FILE <<- EOM
  - job_name: '$NAME'
    metrics_path: /$METRICS_PATH
    dns_sd_configs:
      - names:
        - '$HOST'
        type: 'A'
        port: $PORT

EOM
    done
fi

echo "==prometheus.yml=="
cat $FILE
echo "=================="

echo "Starting Prometheus..."

/bin/prometheus \
    --config.file=/etc/prometheus/prometheus.yml \
    --storage.tsdb.path=/prometheus \
    --web.console.libraries=/usr/share/prometheus/console_libraries \
    --web.console.templates=/usr/share/prometheus/consoles
