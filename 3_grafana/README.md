Grafana
=======

Grafana é um sistema de visualização de dados de analítica e monitoramento<sup>[grafana](https://grafana.com/)</sup>. É possível usá-lo para representar de forma mais elegante as métricas de um sistema capturadas pelo Prometheus (entre outras fontes de dados<sup>[Data Source Guides](https://grafana.com/docs/)</sup>.


Definindo Serviços
------------------

Para subir o serviço do grafana é preciso um **docker-compose.yml**. Neste exemplo tem-se: o **gerador de métricas** para simular um sistema expondo métricas; o **Prometheus** para capturar estas métricas e; o **Grafana** para facilitar a visualização destes dados.

**docker-compose.yml**
```yml
version: '3.3'

services:

  generator:
    image: mfurquim/metrics-generator:1.0.0
    environment:
      - COMPONENT_NAME=testserver
      - COMPONENT_VERSION=1.0.0
    ports:
      - 3000:32865

  grafana:
    image: mfurquim/grafana:5.2.4
    build: ./grafana/
    ports:
      - 4000:3000
    volumes:
      - grafana:/data

  prometheus:
    image: mfurquim/prometheus:v2.7.2
    build: ./prometheus/
    ports:
      - 9090:9090
    environment:
      - STATIC_SCRAPE_TARGETS=generator@generator:32865
    volumes:
      - prometheus:/prometheus

volumes:
  grafana:
  prometheus:
```

Construindo a Imagem
--------------------

O **Dockerfile** para construir a imagem Docker possui poucas linhas. Utilizando a própria imagem do grafana como base, é instalado um plugin como exemplo<sup>[lista plugins](https://grafana.com/grafana/plugins) e adicionado os arquivos definindo a fonte de dados e o local dos dashboards.


```Dockerfile
FROM grafana/grafana:6.1.4

RUN grafana-cli plugins install mtanda-histogram-panel

ADD provisioning /etc/grafana/provisioning
```

**datasource.yml**
```yml
apiVersion: 1

deleteDatasources:
  - name: Prometheus

datasources:
- name: Prometheus
  type: prometheus
  access: proxy
  url: http://prometheus:9090
  isDefault: true
  version: 1
  editable: true
apiVersion: 1
```

**dashboards.yml**
```yml
providers:
- name: 'default'
  orgId: 1
  folder: ''
  type: file
  disableDeletion: false
  editable: true
  options:
    path: /etc/grafana/provisioning/dashboards
```

```
$ docker-compose build
generator uses an image, skipping
Building grafana
[...]
Successfully built 134edb170de6
Successfully tagged mfurquim/grafana:5.2.4
Building prometheus
[...]
Successfully built a2400e997d7c
Successfully tagged mfurquim/prometheus:v2.7.2
```


Rodando o contêiner
-------------------


```
$ docker-compose up
Creating network "3_grafana_default" with the default driver
Creating volume "3_grafana_grafana" with default driver
Creating volume "3_grafana_prometheus" with default driver
Pulling generator (mfurquim/metrics-generator:1.0.0)...
1.0.0: Pulling from mfurquim/metrics-generator
8e402f1a9c57: Pull complete
64ed92a4a25e: Pull complete
Digest: sha256:68d4dc9458032f3009073faa9f2fa4c6c49456e2e808de199c780f9acc3a3dd1
Status: Downloaded newer image for mfurquim/metrics-generator:1.0.0
Creating 3_grafana_generator_1  ... done
Creating 3_grafana_grafana_1    ... done
Creating 3_grafana_prometheus_1 ... done
Attaching to 3_grafana_generator_1, 3_grafana_grafana_1, 3_grafana_prometheus_1
grafana_1     | t=2019-09-12T18:10:20+0000 lvl=info msg="Starting Grafana" logger=server version=6.1.4 commit=fef1733 branch=HEAD compiled=2019-04-16T09:04:07+0000
generator_1   | time="2019-09-12T18:10:20Z" level=info msg="Registering metrics collectors..."
prometheus_1  | Generating prometheus.yml according to ENV variables...
[...]
```

Criando dashboards
------------------

Nesta seção será criado um dashboard com três gráficos por uri, cada um sendo a combinação dentre os três sinais de ouro (tráfego, latência, e erro).

### Adicionando variáveis

Primeira coisa é criar uma variável chamada `uri` para filtrar as _queries_. Ao entrar no grafana através do `localhost:4000`, há dois símbolos de engrenagem. Um à esquerda e um em cima à direita. Para adicionar a variável, clique na engrenagem no canto direito em cima (_dashboard settings_). Em _settings_, a terceira seção está escrito `{x} Variables`. Clique nela e clique em `{x} Add variable` que vai aparecer no centro da tela.

Na tela de adicionar variável, altere os seguintes campos com os valores em itálico: **Name** _uri_; **Data source** _Prometheus_; **Refresh** _On Time Range Change_; **Query** _http\_requests\_throughput_; **Regex** _/.\*uri="(.\*)".\*/_; **Multi-value** _True_. Clique `Add` para adicionar a variável à _dashboard_. Desta forma, aparecerá uma caixa de seleção abaixo do nome do _dashboard_ com as opções de valores das variáveis.

### Criando primeiro painel

Na página principal do _dashboard_, há um painel vazio esperando para ser editado. Clique em _Add Query_ e escreva as _queries_ **A** e **B** como as seguintes.
* A. **query** [http\_requests\_throughput{uri="$uri"}]; **Legend** [Tráfego].
* B. **query** [http\_requests\_latency{uri="$uri"}]; **Legend** [Latência].

























































