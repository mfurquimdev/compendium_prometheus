Counter
=======

Para instrumentar o código adicionando um counter, será usado um programa em linguagem Go. A implementação em outras linguagens é de forma similar. Será usado um exemplo em Java para o exemplo de métrica do histogram.

O conteúdo do arquivo `src/main.go` será mostrado em partes para o melhor entendimento do código. Para executá-lo é preciso instalar o [Docker](https://www.docker.com/get-started) para construir a imagem já contendo o programa compilado, ou instalar [golang](https://golang.org/dl/) localmente e baixar a dependência

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)
```

A primeira parte do programa em Go é indicar a qual pacote este arquivo pertence. Ao rodar um programa em Golang, o Go procura por um pacote e uma função com nome _main_ para executar primeiro. Caso não haja, é retornado um erro. As próximas linhas são de importação de dependência: a lib _log_ é para escrever na saída padrão caso haja um erro; a _net/http_ serve para abrir um serviço web na porta 8080; a _time_ será usada apenas para esperar um segundo antes de incrementar a variável; _prometheus_ e _promhttp_ implementam as métricas e a forma de expô-las.

```go
func main() {
  // Creates new/empty Registry
  upTimeReg := prometheus.NewRegistry()
```

Após importar todas as bibliotecas necessárias, é criado um registrador. **Registrador** é uma estrutura que armazena as métricas e pode ser exposto em um _endpoint_ específico. Nas linhas abaixo, a métrica será declarada e registrada. No final do arquivo, o `upTimeReg` será exposto na porta `8080` e caminho `/metrics-counter`.

```go
  // A counter metric for how long (in seconds) the server is up
  var upTime = prometheus.NewCounter(
    prometheus.CounterOpts{
      Name: "uptime_seconds_total",
      Help: "Time in seconds the service is up",
    },
  )
```

Nas linhas acima, é declarado a variável que será utilizado para a métrica. O exemplo utiliza apenas um contador, passando os parâmetros de nome e uma pequena frase de ajuda para auxiliar o seu entendimento. Um **contador** é um tipo de métrica que pode apenas ter seu valor incrementado em números positivos.

```go

  // Register metric in the Registry
  upTimeReg.MustRegister(upTime)
```

A função acima registra uma métrica em um registrador, para enfim expô-la.

```go
	// Increment variable each second inside a goroutine (parallel)
	go func() {
		for {
			time.Sleep(time.Second)
			upTime.Inc()
		}
	}()
```

O contador que está sendo usado será incrementado uma vez por segundo. A palavra-chave `go` executa a função em paralelo. A função anônima, por sua vez, espera um segundo e incrementa a variável `upTime` em um loop infinito. Esta não é a melhor maneira de verificar o tempo que um servidor está rodando, mas é um exemplo simples o suficiente para entender a biblioteca do _prometheus_.


```go
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
  http.Handle("/metrics-counter", promhttp.HandlerFor(upTimeReg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Por fim, um servidor http é aberto na porta `8080`, expondo as métricas no _endpoint_ `/metrics-counter`. Ao rodar o programa, seja localmente ou com Docker, a métrica `uptime_seconds_total` será exposta. Para acessar as métricas, digite `http://localhost:8080/metrics-counter` no browser. O conteúdo será parecido com o seguinte: 

```
# HELP uptime_seconds_total Time in seconds the service is up
# TYPE uptime_seconds_total counter
uptime_seconds_total 4
```

Docker e docker-compose
-----------------------

Para executar o programa em um Docker, é preciso de um **Dockerfile** e, preferencialmente, um **docker-compose.yml** também. O **Dockerfile** para definir como a imagem deve ser construída, e o **docker-compose.yml** para definir como ela deve ser executada.

**Dockerfile**
```Dockerfile
FROM golang:alpine

RUN apk update && apk add git 

COPY src/main.go $GOPATH/src/promster/main.go

RUN CGO_ENABLED=0 GOOS=linux go get -v promster

EXPOSE 8080/tcp

ENTRYPOINT [ "/go/bin/promster" ]
```

O arquivo de descrição da imagem importa uma imagem para usar como base. É instalado o `git` para baixar as dependências, o arquivo `main.go` copiado para dentro da imagem e, depois, compilado. A porta 8080 é exposta e o programa para iniciar quando o contêiner for rodado é definido em `ENTRYPOINT`.


**docker-compose.yml**
```yml
version: '3.5'

services:
  metric:
    build: .
    image: mfurquim/counter:v1.0.0
    ports:
      - 8080:8080
```

A porta 8080 do contêiner é mapeada para a porta 8080 do host. Para executar, use o seguinte comando: `$ docker-compose up --build`

Golang local
------------

Para executar o programa localmente, instale o **Go**, baixe as dependências com `$ go get github.com/prometheus/client_golang/prometheus`, e rode o programa com `$ go run src/main.go`. Ou compile-o com `$ go build src/main.go` e execute com `$ ./main`.
