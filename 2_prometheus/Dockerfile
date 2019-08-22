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
