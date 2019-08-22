#!/bin/bash

# WARNING: This code has ugly implementations because Prometheus base image doesn't have bash, just sh!

echo "Generating prometheus.yml according to ENV variables..."
echo "Informed variables:"
echo "1 - SCRAPE_INTERVAL= $SCRAPE_INTERVAL"
echo "2 - EVALUATION_INTERVAL=$EVALUATION_INTERVAL"
echo "3 - SCRAPE_TIMEOUT=$SCRAPE_TIMEOUT"
echo "4 - TSDB_RETENTION=$TSDB_RETENTION"
echo "5 - REGION=$REGION"
echo "6 - SERVICE=$SERVICE"
echo "7 - PROMETHEUS_NAME=$PROMETHEUS_NAME"

# SANITY CHECK
if [ "$SCRAPE_INTERVAL" == "" ]; then
  echo "SCRAPE_INTERVAL ENV is required" 1>&2
  exit 1
fi

if [ "$EVALUATION_INTERVAL" == "" ]; then
  echo "EVALUATINO_INTERVAL ENV is required" 1>&2
  exit 1
fi

if [ "$SCRAPE_TIMEOUT" == "" ]; then
  echo "SCRAPE_TIMEOUT ENV is required" 1>&2
  exit 1
fi

if [ "$TSDB_RETENTION" == "" ]; then
  echo "TSDB_RETENTION ENV is required" 1>&2
  exit 1
fi

FILE=/etc/prometheus/prometheus.yml

#### GLOBAL DEFINITIONS ####
cat > $FILE <<- EOM
global:
  scrape_interval: $SCRAPE_INTERVAL
  evaluation_interval: $EVALUATION_INTERVAL
  scrape_timeout: $SCRAPE_TIMEOUT

EOM

RULES=""
NEWLINE=$'\n'
for file in /etc/prometheus/*.yml; do
    FILENAME="$(expr "$file" : '/etc/prometheus/\(.*\)')"
    if [ ! "$FILENAME" == "prometheus.yml" ]; then
        RULES="${RULES}${NEWLINE}  - ${FILENAME}"
        sed -i -e 's/$REGION/'"$REGION"'/g' '/etc/prometheus/'$FILENAME
        sed -i -e 's/$PROMETHEUS_NAME/'"$PROMETHEUS_NAME"'/g' '/etc/prometheus/'$FILENAME
        sed -i -e 's/$EVALUATION_INTERVAL/'"$EVALUATION_INTERVAL"'/g' '/etc/prometheus/'$FILENAME
    fi
done

cat >> $FILE <<- EOM
rule_files: $RULES

EOM


cat >> $FILE <<- EOM
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
    - targets: ['localhost:9090']

EOM

for MP in $(echo ${METRICS_PATHS} | tr "," "\n" )
do

cat >> $FILE <<- EOM
  - job_name: "metrics_generator$MP"
    metrics_path: $MP
    scheme: $SCHEME
    tls_config:
      insecure_skip_verify: true
    file_sd_configs:
      - files:
        - targets.json
EOM

done

echo "==prometheus.yml=="
cat $FILE
echo "=================="

echo "Starting Prometheus..."

/bin/prometheus \
    --config.file=/etc/prometheus/prometheus.yml \
    --storage.tsdb.path=/prometheus \
    --storage.tsdb.retention="$TSDB_RETENTION" \
    --web.console.libraries=/usr/share/prometheus/console_libraries \
    --web.console.templates=/usr/share/prometheus/consoles
