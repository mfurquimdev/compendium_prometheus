FROM golang:alpine

RUN apk update && apk add git 

COPY src/main.go $GOPATH/src/promster/main.go

RUN CGO_ENABLED=0 GOOS=linux go get -v promster

EXPOSE 8080/tcp

ENTRYPOINT [ "/go/bin/promster" ]

