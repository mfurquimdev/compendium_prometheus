Histogram
=========

Para este segundo exemplo, será usado um código Java para expor a métrica de _histogram_. Será mostrado apenas as partes importantes para importar as bibliotecas, declarar variáveis, alterá-la, e expô-la.

Para gerar o pacote, será usado o maven. O maven precisa de um arquivo **pom.xml** que contém informações sobre o projeto e suas dependências. As dependências importantes para este exemplo são:

**pom.xml**
```xml
<dependencies>
    <dependency>
        <groupId>javax.servlet</groupId>
        <artifactId>javax.servlet-api</artifactId>
        <version>4.0.1</version>
    </dependency>

    <dependency>
        <groupId>io.prometheus</groupId>
        <artifactId>simpleclient</artifactId>
        <version>0.6.0</version>
    </dependency>

    <dependency>
        <groupId>io.prometheus</groupId>
        <artifactId>simpleclient_servlet</artifactId>
        <version>0.6.0</version>
    </dependency>
</dependencies>
```

No **TestServlet.java**, é importado o contador, histograma e o exportador das métricas do prometheus.

```java
import io.prometheus.client.Counter;
import io.prometheus.client.Histogram;
import io.prometheus.client.exporter.MetricsServlet;
```

Depois, é declarado as variáveis utilizando o build do tipo de métrica que será utilizado. O endpoint para este teste é `/test`. Isso quer dizer que cada vez que acessar o `{{url}}/test`, a variável será incrementada de acordo.

```java
@WebServlet(name = "TestServlet", urlPatterns = "/test")
public class TestServlet extends HttpServlet {
    static final Counter requests = Counter.build()
            .name("requests_total").help("Total número de requisições.").register();

    static final Histogram histRand = Histogram.build()
            .buckets(0.1, 0.25, 0.5, 0.75, 0.9, 1)
            .name("requests_random_numbers").help("Random number generated").register();

```

Tanto o `Counter` quanto o `Histogram` precisam de nomes e de um texto de ajuda, assim como no exemplo passado em Golang. Depois de criá-las, é preciso registrar com o `.register()`. O `Histogram` ainda possui um parâmetro a mais opcional que são os "baldes". Estes _buckets_ nada mais são que contadores dos valores. Foram definidos alguns valores entre 0 e 1. Quando um valor abaixo de ou igual a (`le`) 1.0 é observado, o _bucket_ referente a 1.0 vai incrementar. Caso um valor `le=0.9` seja observado, então os baldes tanto de 1.0 quanto de 0.9 serão incrementados. Isso por que o valor que é menor do que 0.9, também é menor do que 1.0. Dessa forma, um valor menor que 0.1 (o menor balde) irá incrementar todos os _buckets_. Mais abaixo será mostrado um possível resultado no `/metrics`.

Quando o `{{url}}/test` é acessado, a função `doGet()` é executada, de forma que o `requests` é incrementado e `histRand` é adicionado um valor aleatório entre 0 e 1.

```java
    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
        requests.inc();
        histRand.observe(Math.random());
```

Para expor o endpoint `/metrics` foi utilizado a anotação `@WebServlet` de uma classe que extende do MetricsServlet.

**MServlet.java**
```java
package br.com.furquim.servlets;

import io.prometheus.client.exporter.MetricsServlet;

import javax.servlet.annotation.WebServlet;

@WebServlet(name = "MServlet", urlPatterns = "/metrics")
public class MServlet extends MetricsServlet {
}
```



Para gerar o arquivo `war` é preciso executar o seguinte comando do maven:

```
$ mvn package
[INFO] Scanning for projects...
[WARNING]
[WARNING] Some problems were encountered while building the effective model for br.com.furquim:servlet-test:war:1.0-SNAPSHOT
[WARNING] 'build.plugins.plugin.version' for org.apache.maven.plugins:maven-war-plugin is missing. @ line 12, column 21
[WARNING]
[WARNING] It is highly recommended to fix these problems because they threaten the stability of your build.
[WARNING]
[WARNING] For this reason, future Maven versions might no longer support building such malformed projects.
[WARNING]
[INFO]
[INFO] --------------------< br.com.furquim:servlet-test >---------------------
[INFO] Building servlet-test 1.0-SNAPSHOT
[INFO] --------------------------------[ war ]---------------------------------
[INFO]
[INFO] --- maven-resources-plugin:2.6:resources (default-resources) @ servlet-test ---
[WARNING] Using platform encoding (UTF-8 actually) to copy filtered resources, i.e. build is platform dependent!
[INFO] Copying 0 resource
[INFO]
[INFO] --- maven-compiler-plugin:3.1:compile (default-compile) @ servlet-test ---
[INFO] Changes detected - recompiling the module!
[WARNING] File encoding has not been set, using platform encoding UTF-8, i.e. build is platform dependent!
[INFO] Compiling 2 source files to /Users/mfurquim/IBM/compendio_prometheus/1_metrics/3_histogram/target/classes
[INFO]
[INFO] --- maven-resources-plugin:2.6:testResources (default-testResources) @ servlet-test ---
[WARNING] Using platform encoding (UTF-8 actually) to copy filtered resources, i.e. build is platform dependent!
[INFO] skip non existing resourceDirectory /Users/mfurquim/IBM/compendio_prometheus/1_metrics/3_histogram/src/test/resources
[INFO]
[INFO] --- maven-compiler-plugin:3.1:testCompile (default-testCompile) @ servlet-test ---
[INFO] Nothing to compile - all classes are up to date
[INFO]
[INFO] --- maven-surefire-plugin:2.12.4:test (default-test) @ servlet-test ---
[INFO] No tests to run.
[INFO]
[INFO] --- maven-war-plugin:2.2:war (default-war) @ servlet-test ---
[INFO] Packaging webapp
[INFO] Assembling webapp [servlet-test] in [/Users/mfurquim/IBM/compendio_prometheus/1_metrics/3_histogram/target/servlet-test-1.0-SNAPSHOT]
[INFO] Processing war project
[INFO] Webapp assembled in [62 msecs]
[INFO] Building war: /Users/mfurquim/IBM/compendio_prometheus/1_metrics/3_histogram/target/servlet-test-1.0-SNAPSHOT.war
[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ------------------------------------------------------------------------
[INFO] Total time:  2.561 s
[INFO] Finished at: 2019-08-07T17:18:29-03:00
[INFO] ------------------------------------------------------------------------
```

O arquivo `war` será criado no diretório `target`. O `docker-compose` cuida de construir a imagem do `tomcat` contendo o `war` e executar o contêiner.

**docker-compose.yml**
```
version: '3.5'

services:
  tomcat:
    image: mfurquim/metrics-test:v0.0.1
    build: .
    ports:
      - "8888:8080"
```

No `ports`, a porta `8080` do contêiner é mapeada para a porta `8888` do host.

**Dockerfile**
```
FROM tomcat:8.0

COPY target/servlet-test-1.0-SNAPSHOT.war /usr/local/tomcat/webapps/
COPY tomcat-users.xml  $CATALINA_HOME/conf/
```

O `Dockerfile` copia o `war`para o `webapps`, que o Tomcat vai olhar para fazer o deploy automático. O outro arquivo adicionado para dentro da imagem do Tomcat é um arquivo de usuários para caso queira fazer login e gerenciar os deploys.

**tomcat-users.xml**
```
<?xml version='1.0' encoding='utf-8'?>
<tomcat-users>
  <role rolename="manager"/>
  <role rolename="admin"/>
  <user username="tomcat" password="tomcat" roles="admin,manager, manager-gui"/>
</tomcat-users>
```

Execute o seguinte comando para construir a imagem e executar o contêiner:

```
$ docker-compose up --build
Building tomcat
Step 1/3 : FROM tomcat:8.0
 ---> ef6a7c98d192
Step 2/3 : COPY target/servlet-test-1.0-SNAPSHOT.war /usr/local/tomcat/webapps/
 ---> e54b5b00b6d1
Step 3/3 : COPY tomcat-users.xml  $CATALINA_HOME/conf/
 ---> 87ad7a6d03b0

Successfully built 87ad7a6d03b0
Successfully tagged mfurquim/metrics-test:v0.0.1
Creating 3_histogram_tomcat_1 ... done
Attaching to 3_histogram_tomcat_1
tomcat_1  | 07-Aug-2019 20:20:26.140 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Server version:        Apache Tomcat/8.0.53
tomcat_1  | 07-Aug-2019 20:20:26.147 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Server built:          Jun 29 2018 14:42:45 UTC
tomcat_1  | 07-Aug-2019 20:20:26.147 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Server number:         8.0.53.0
tomcat_1  | 07-Aug-2019 20:20:26.147 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log OS Name:               Linux
tomcat_1  | 07-Aug-2019 20:20:26.147 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log OS Version:            4.9.184-linuxkit
tomcat_1  | 07-Aug-2019 20:20:26.147 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Architecture:          amd64
tomcat_1  | 07-Aug-2019 20:20:26.148 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Java Home:             /usr/lib/jvm/java-7-openjdk-amd64/jre
tomcat_1  | 07-Aug-2019 20:20:26.148 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log JVM Version:           1.7.0_181-b01
tomcat_1  | 07-Aug-2019 20:20:26.148 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log JVM Vendor:            Oracle Corporation
tomcat_1  | 07-Aug-2019 20:20:26.148 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log CATALINA_BASE:         /usr/local/tomcat
tomcat_1  | 07-Aug-2019 20:20:26.148 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log CATALINA_HOME:         /usr/local/tomcat
tomcat_1  | 07-Aug-2019 20:20:26.149 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Command line argument: -Djava.util.logging.config.file=/usr/local/tomcat/conf/logging.properties
tomcat_1  | 07-Aug-2019 20:20:26.149 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Command line argument: -Djava.util.logging.manager=org.apache.juli.ClassLoaderLogManager
tomcat_1  | 07-Aug-2019 20:20:26.149 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Command line argument: -Djdk.tls.ephemeralDHKeySize=2048
tomcat_1  | 07-Aug-2019 20:20:26.149 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Command line argument: -Djava.protocol.handler.pkgs=org.apache.catalina.webresources
tomcat_1  | 07-Aug-2019 20:20:26.150 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Command line argument: -Dignore.endorsed.dirs=
tomcat_1  | 07-Aug-2019 20:20:26.150 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Command line argument: -Dcatalina.base=/usr/local/tomcat
tomcat_1  | 07-Aug-2019 20:20:26.150 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Command line argument: -Dcatalina.home=/usr/local/tomcat
tomcat_1  | 07-Aug-2019 20:20:26.150 INFO [main] org.apache.catalina.startup.VersionLoggerListener.log Command line argument: -Djava.io.tmpdir=/usr/local/tomcat/temp
tomcat_1  | 07-Aug-2019 20:20:26.150 INFO [main] org.apache.catalina.core.AprLifecycleListener.lifecycleEvent Loaded APR based Apache Tomcat Native library 1.2.17 using APR version 1.5.1.
tomcat_1  | 07-Aug-2019 20:20:26.150 INFO [main] org.apache.catalina.core.AprLifecycleListener.lifecycleEvent APR capabilities: IPv6 [true], sendfile [true], accept filters [false], random [true].
tomcat_1  | 07-Aug-2019 20:20:26.154 INFO [main] org.apache.catalina.core.AprLifecycleListener.initializeSSL OpenSSL successfully initialized (OpenSSL 1.1.0f  25 May 2017)
tomcat_1  | 07-Aug-2019 20:20:26.238 INFO [main] org.apache.coyote.AbstractProtocol.init Initializing ProtocolHandler ["http-apr-8080"]
tomcat_1  | 07-Aug-2019 20:20:26.249 INFO [main] org.apache.coyote.AbstractProtocol.init Initializing ProtocolHandler ["ajp-apr-8009"]
tomcat_1  | 07-Aug-2019 20:20:26.251 INFO [main] org.apache.catalina.startup.Catalina.load Initialization processed in 641 ms
tomcat_1  | 07-Aug-2019 20:20:26.290 INFO [main] org.apache.catalina.core.StandardService.startInternal Starting service Catalina
tomcat_1  | 07-Aug-2019 20:20:26.290 INFO [main] org.apache.catalina.core.StandardEngine.startInternal Starting Servlet Engine: Apache Tomcat/8.0.53
tomcat_1  | 07-Aug-2019 20:20:26.307 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployWAR Deploying web application archive /usr/local/tomcat/webapps/servlet-test-1.0-SNAPSHOT.war
tomcat_1  | 07-Aug-2019 20:20:26.928 INFO [localhost-startStop-1] org.apache.jasper.servlet.TldScanner.scanJars At least one JAR was scanned for TLDs yet contained no TLDs. Enable debug logging for this logger for a complete list of JARs that were scanned but no TLDs were found in them. Skipping unneeded JARs during scanning can improve startup time and JSP compilation time.
tomcat_1  | 07-Aug-2019 20:20:26.977 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployWAR Deployment of web application archive /usr/local/tomcat/webapps/servlet-test-1.0-SNAPSHOT.war has finished in 669 ms
tomcat_1  | 07-Aug-2019 20:20:26.978 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployDirectory Deploying web application directory /usr/local/tomcat/webapps/examples
tomcat_1  | 07-Aug-2019 20:20:27.385 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployDirectory Deployment of web application directory /usr/local/tomcat/webapps/examples has finished in 407 ms
tomcat_1  | 07-Aug-2019 20:20:27.385 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployDirectory Deploying web application directory /usr/local/tomcat/webapps/host-manager
tomcat_1  | 07-Aug-2019 20:20:27.433 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployDirectory Deployment of web application directory /usr/local/tomcat/webapps/host-manager has finished in 47 ms
tomcat_1  | 07-Aug-2019 20:20:27.433 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployDirectory Deploying web application directory /usr/local/tomcat/webapps/docs
tomcat_1  | 07-Aug-2019 20:20:27.466 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployDirectory Deployment of web application directory /usr/local/tomcat/webapps/docs has finished in 33 ms
tomcat_1  | 07-Aug-2019 20:20:27.466 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployDirectory Deploying web application directory /usr/local/tomcat/webapps/ROOT
tomcat_1  | 07-Aug-2019 20:20:27.510 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployDirectory Deployment of web application directory /usr/local/tomcat/webapps/ROOT has finished in 44 ms
tomcat_1  | 07-Aug-2019 20:20:27.511 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployDirectory Deploying web application directory /usr/local/tomcat/webapps/manager
tomcat_1  | 07-Aug-2019 20:20:27.541 INFO [localhost-startStop-1] org.apache.catalina.startup.HostConfig.deployDirectory Deployment of web application directory /usr/local/tomcat/webapps/manager has finished in 31 ms
tomcat_1  | 07-Aug-2019 20:20:27.550 INFO [main] org.apache.coyote.AbstractProtocol.start Starting ProtocolHandler ["http-apr-8080"]
tomcat_1  | 07-Aug-2019 20:20:27.580 INFO [main] org.apache.coyote.AbstractProtocol.start Starting ProtocolHandler ["ajp-apr-8009"]
tomcat_1  | 07-Aug-2019 20:20:27.586 INFO [main] org.apache.catalina.startup.Catalina.start Server startup in 1335 ms
```

Agora que o servidor de aplicações Tomcat está rodando com o `servlet-test`, acesse a url `http://localhost:8888/servlet-test-1.0-SNAPSHOT/test` algumas vezes para gerar as métricas. Um possível resultado depois de acessar 49 vezes é o seguinte:
```
Número de requisições: 49.0
Histogram Random:
	Name: requests_random_numbers_bucket LabelNames: [le] labelValues: [0.1] Value: 8.0 TimestampMs: null
	Name: requests_random_numbers_bucket LabelNames: [le] labelValues: [0.25] Value: 14.0 TimestampMs: null
	Name: requests_random_numbers_bucket LabelNames: [le] labelValues: [0.5] Value: 26.0 TimestampMs: null
	Name: requests_random_numbers_bucket LabelNames: [le] labelValues: [0.75] Value: 39.0 TimestampMs: null
	Name: requests_random_numbers_bucket LabelNames: [le] labelValues: [0.9] Value: 46.0 TimestampMs: null
	Name: requests_random_numbers_bucket LabelNames: [le] labelValues: [1.0] Value: 49.0 TimestampMs: null
	Name: requests_random_numbers_bucket LabelNames: [le] labelValues: [+Inf] Value: 49.0 TimestampMs: null
	Name: requests_random_numbers_count LabelNames: [] labelValues: [] Value: 49.0 TimestampMs: null
	Name: requests_random_numbers_sum LabelNames: [] labelValues: [] Value: 22.38473968542088 TimestampMs: null
```

O Endpoint que o Prometheus faz scrap das métricas é `http://localhost:8888/servlet-test-1.0-SNAPSHOT/metrics`. O mesmo resultado das métricas quando acessado o `/test` 49 vezes, mas no `/metrics`, é o seguinte:
```
# HELP requests_total Total número de requisições.
# TYPE requests_total counter
requests_total 49.0
# HELP requests_random_numbers Random number generated
# TYPE requests_random_numbers histogram
requests_random_numbers_bucket{le="0.1",} 8.0
requests_random_numbers_bucket{le="0.25",} 14.0
requests_random_numbers_bucket{le="0.5",} 26.0
requests_random_numbers_bucket{le="0.75",} 39.0
requests_random_numbers_bucket{le="0.9",} 46.0
requests_random_numbers_bucket{le="1.0",} 49.0
requests_random_numbers_bucket{le="+Inf",} 49.0
requests_random_numbers_count 49.0
requests_random_numbers_sum 22.38473968542088
```



