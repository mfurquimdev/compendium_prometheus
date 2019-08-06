Compêndio Prometheus + Grafana
==============================

Docker
--------

Todos os exemplos utilizarão `Docker` e `docker-compose`. Este conhecimento é essencial e será explicado na sessão [0\_docker].

1 /Metrics
----------

Instrumentar códgio em Java para expor /metrics


2 Prometheus
------------

Subir um prometheus para fazer scrape e mostrá-lo no /graph


Federação: conceito e como fazer para passar info no /federate


3 Grafana
---------

Arquivo de configuração para adicionar o prometheus como fonte dos dados


4 Promster
----------

Subir um conjunto de serviços para o promster rodar com 2 níveis
1. Subir ETCD para escrever, ler, e apagar uma chave (explicar `ETCDCTL_API=3` e como dar `export`)
2. Subir um `etcd_registrar` para registrar um IP e mantê-lo ativo, matar o serviço para ver a chave sumindo do ETCD (utilizar `watch etcdctl get --prefix /`)
3. Subir o `promster` level 1 e level 2
4. Configurar o grana
