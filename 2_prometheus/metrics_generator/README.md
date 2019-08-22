# metrics-generator-tabajara

Randomly generates Prometheus metrics for simulating default /metrics endpoints accross servers

### Usage

docker-compose.yml
```
version: '3.3'

services:

  generator:
    image: labbsr0x/metrics-generator-tabajara
    build: .
    environment:
      - COMPONENT_NAME=testserver
      - COMPONENT_VERSION=1.0.0
    ports:
      - 3000:3000
```

curl http://localhost:3000/metrics

### Cause surgery accidents!

* If you want to change any resource metric abnormally to test alert rules, for example, do

```
curl -X POST http://localhost:3000/surgery-accident
{
	"resource": "transaction-0100",
	"type": "latency",
	"value": "100"
}
```

* After this, transaction-0100 median latency will be ~100x the original

* 'resource' accepts regular expressions for matching N resource names (for example, 'transaction-.*')

* To remove all accidents, send a DELETE to the same endpoint

```
curl -X DELETE http://localhost:3000/surgery-accident
```
