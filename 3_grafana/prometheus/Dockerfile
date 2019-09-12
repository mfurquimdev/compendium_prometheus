FROM debian:stretch

RUN apt-get update && apt-get install -y --no-install-recommends \
    gettext \
    && rm -rf /var/lib/apt/lists/*

FROM prom/prometheus:v2.7.2

COPY --from=0 /usr/bin/envsubst /usr/bin/envsubst
RUN ls -al /usr/bin/envsubst

ENV SCRAPE_INTERVAL 15s
ENV EVALUATION_INTERVAL 15s
ENV SCRAPE_TIMEOUT 10s

USER root
ADD rules /etc/prometheus
ADD startup.sh /
ADD build.sh /

RUN chmod +x /build.sh
RUN sh /build.sh /etc/prometheus/

ENTRYPOINT [ "/bin/sh" ]
CMD [ "/startup.sh" ]
