Gauge
=====

Este exemplo do tipo de métrica _Gauge_ também será escrito em Go. O conteúdo do arquivo `src/main.go` será mostrado em partes para o melhor entendimento do código. Para executá-lo é preciso instalar o [Docker](https://www.docker.com/get-started) para construir a imagem já contendo o programa compilado, ou instalar [golang](https://golang.org/dl/) localmente e baixar a dependência

```go
package main

import (
	"log"
  "math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)
```

A primeira parte do programa em Go é indicar a qual pacote este arquivo pertence. Ao rodar um programa em Golang, o Go procura por um pacote e uma função com nome _main_ para executar primeiro. Caso não haja, é retornado um erro. As próximas linhas são de importação de dependência: a lib _log_ é para escrever na saída padrão caso haja um erro; a _math/rand_ é utilizada para gerar números aleatórios; a _net/http_ serve para abrir um serviço web na porta 8080; a _time_ será usada apenas para esperar um tempo, em segundos, para finalizar a função e alterar a variável da métrica; _prometheus_ e _promhttp_ implementam as métricas e a forma de expô-las.

```go
func main() {
  // Creates new/empty Registry
  funcDurationReg := prometheus.NewRegistry()
```

Após importar todas as bibliotecas necessárias, é criado um registrador vazio para a métrica.

```go
  // A counter metric for how long (in seconds) the server is up
  var funcDuration = prometheus.NewGauge(
    prometheus.GaugeOpts{
      Name: "function_duration_seconds",
      Help: "Time in seconds the most recently run of a function has taken to complete",
    },
  )
```

Nas linhas acima, é declarado a variável que será utilizado para a métrica. O exemplo utiliza apenas um _gauge_, passando os parâmetros de nome e uma pequena frase de ajuda para auxiliar o seu entendimento. Um **gauge** é um tipo de métrica que pode ter seu valor alterado arbitrariamente, tanto incrementado quanto decrementando.

```go
  // Register metric in the Registry
  funcDurationReg.MustRegister(funcDuration)
```

A função acima registra uma métrica em um registrador, para enfim expô-la.

```go
  // Execute a function in a goroutine to alter the variable funcDuration
  // based on a random time waiting
  go func() {
    for {
      go func() {
        timer := prometheus.NewTimer(prometheus.ObserverFunc(funcDuration.Set))
        defer timer.ObserveDuration()
        time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
      }()
      time.Sleep(5 * time.Second)
    }
  }()
```

A função acima dispara uma rotina go a cada cinco segundos para executar uma função que, ao final, registre o tempo que levou para terminar. A palavra-chave `defer` executa o comando seguinte apenas quando o escopo da função terminar. Quando a função terminar, o tempo gasto será observado e o valor atribuído à variável `funcDuration`. A função espera um tempo, em segundos, aleatório no intervalo de [0,1).


```go
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
  http.Handle("/metrics-gauge", promhttp.HandlerFor(funcDurationReg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Por fim, um servidor http é aberto na porta `8080`, expondo as métricas no _endpoint_ `/metrics-gauge`. Ao rodar o programa, seja localmente ou com Docker, a métrica `function_duration_seconds` será exposta. Para acessar as métricas, digite `http://localhost:8080/metrics-gauge` no browser. O conteúdo será parecido com o seguinte: 

```
# HELP function_duration_seconds Time in seconds the most recently run of a function has taken to complete
# TYPE function_duration_seconds gauge
function_duration_seconds 8.003958078
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
    image: mfurquim/gauge:v1.0.0
    ports:
      - 8080:8080
```

A porta 8080 do contêiner é mapeada para a porta 8080 do host. Para executar, use o seguinte comando: `$ docker-compose up --build`

Golang local
------------

Para executar o programa localmente, instale o **Go**, baixe as dependências com `$ go get github.com/prometheus/client_golang/prometheus`, e rode o programa com `$ go run src/main.go`. Ou compile-o com `$ go build src/main.go` e execute com `$ ./main`.

