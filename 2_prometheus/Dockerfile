FROM prom/prometheus:v2.7.2

#### ARGS #####

ARG scrapeinterval
ARG evaluationinterval
ARG scrapetimeout
ARG tsdbretention
# targetsfile defines the file path of the targets a prometheus instance of this image will be hitting
ARG targetsfile
# rulespath defines the rules folder a prometheus instance of this image will be using
ARG rulespath
# startupfile defines the file path of the startup script
ARG startupfile
ARG prometheusname
ARG region
ARG service

#### ENVS ####

# global envs
ENV SCRAPE_INTERVAL ${scrapeinterval}
ENV EVALUATION_INTERVAL ${evaluationinterval}
ENV SCRAPE_TIMEOUT ${scrapetimeout}
# tsdb storage envs
ENV TSDB_RETENTION ${tsdbretention}
# REGION defines the deployment region (blue|green)
ENV REGION ${region}
# SERVICE defines which kind of service should a l1 prometheus consume (global|release)
ENV SERVICE ${service}
# PROMETHEUS_NAME identifies the prometheus instance
ENV PROMETHEUS_NAME ${prometheusname}

#### CONFIG ####

USER root

ADD $targetsfile /etc/prometheus/targets.json
ADD $rulespath /etc/prometheus/
ADD $startupfile /

ADD prometheus.yml /etc/prometheus/
ADD build.sh /

RUN chmod -R 777 /etc/prometheus/
RUN chmod -R 777 /startup.sh
RUN chmod +x /build.sh

RUN sh /build.sh /etc/prometheus/

ENTRYPOINT [ "/bin/sh" ]
CMD [ "/startup.sh" ]