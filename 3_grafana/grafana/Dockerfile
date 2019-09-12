FROM grafana/grafana:6.1.4

RUN grafana-cli plugins install mtanda-histogram-panel

ADD provisioning /etc/grafana/provisioning
