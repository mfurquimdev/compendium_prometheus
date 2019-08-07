Counter
=======

Para instrumentar o código adicionando um counter, será usado um programa em linguagem Go. A implementação em outras linguagens é de forma similar. Será usado um exemplo em Java para o exemplo de métrica do histogram.

O conteúdo do arquivo `src/main.go` será mostrado em partes para o melhor entendimento do código. Para executá-lo é preciso instalar [golang](https://golang.org/dl/) e baixar a dependência com `$ go get github.com/prometheus/client_golang/prometheus`.

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

A primeira parte do programa em Go é indicar a qual pacote este arquivo pertence. Ao rodar um programa em Golang, o Go procura por um pacote e uma função _main_ para executar primeiro. Caso não haja, é retornado um erro. As próximas linhas são de importação de dependência: a lib _log_ é para escrever em um arquivo ou saída padrão caso haja um erro; a _net/http_ serve para abrir um serviço web na porta 8080; a _time_ será usada apenas para esperar um segundo antes de incrementar a variável; _prometheus_ e _promhttp_ implementam as métricas e a forma de expor no /metrics.

```go
var timeUp = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "uptime_seconds_total",
		Help: "Time in seconds the service is up",
	},
)
```

Após importar todas as bibliotecas necessárias, é declarado a variável que será utilizado para a métrica. O exemplo utiliza apenas um contador, passando os parâmetros de nome e uma pequena ajuda para ajudar a entendê-la melhor.

```go
func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(timeUp)
}
```

Para expor as métricas, é preciso registrá-las. A função acima faz exatamente isso.

```go
func main() {
	// Increment variable each second inside a goroutine (parallel)
	go func() {
		for {
			time.Sleep(time.Second)
			timeUp.Inc()
		}
	}()
```

O contador que está sendo usado incrementará uma vez por segundo. A palavra-chave `go` executa a função em paralelo. A função anônima, por sua vez, espera um segundo e incrementa a variável `timeUp` em um loop infinito.

```go
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Por fim, um servidor http é aberto na porta `8080`, expondo as métricas no endpoint `/metrics`. Ao rodar o programa com `$ go run main.go`, tem-se diversas métricas do próprio _go_ e do _promhttp_, é preciso procurar um a métrica `uptime_seconds_total` dentre elas. Para acessar as métricas, digite `http://localhost:8080/metrics` no browser.

```
go_*
promhttp_*
# HELP uptime_seconds_total Time in seconds the service is up
# TYPE uptime_seconds_total counter
uptime_seconds_total 19
```
