FROM golang:1.10 AS BUILD

#doing dependency build separated from source build optimizes time for developer, but is not required
#install external dependencies first
ADD /main.go $GOPATH/src/metrics-generator-tabajara/main.go
RUN go get -v metrics-generator-tabajara

#now build source code
ADD metrics-generator-tabajara $GOPATH/src/metrics-generator-tabajara
RUN go get -v metrics-generator-tabajara

FROM flaviostutz/etcd-registrar AS REGISTRAR

FROM golang:1.10 AS IMAGE

EXPOSE 3000

ENV SERVER_NAME ''
ENV COMPONENT_NAME 'testserver'
ENV COMPONENT_VERSION '1.0.0'
ENV ACCIDENT_RESOURCE ''
ENV ACCIDENT_TYPE ''
ENV ACCIDENT_RATIO 1

ENV REGISTRY_ETCD_URL ""
ENV REGISTRY_ETCD_BASE "/services"
ENV REGISTRY_SERVICE "generator"
ENV REGISTRY_TTL "60"

COPY --from=BUILD /go/bin/* /bin/
COPY --from=REGISTRAR /bin/etcd-registrar /bin/etcd-registrar
ADD startup.sh /

CMD [ "/startup.sh" ]
