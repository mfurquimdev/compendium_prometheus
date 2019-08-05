Docker
======

Criação de imagem Docker
------------------------

Para criar uma imagem docker, basta executar o comando `docker build` no diretório em que se encontra o arquivo `Dockerfile`.

**Dockerfile** é o arquivo contendo as instruções de como construir a imagem Docker.

O script que será executado ao subir o contêiner se chama `NossoScript.sh` e se encontra dentro do diretório `./0_DockerImage`. O conteúdo dos arquivos `Dockerfile` e `NossoScript.sh` é o seguinte:

**NossoScript.sh**
```
#!/bin/sh

echo "Olá do nosso script"
```

**Dockerfile**
```
FROM alpine

COPY NossoScript.sh /
RUN chmod a+x /NossoScript.sh 

CMD ["/NossoScript.sh"]
```

O arquivo `NossoScript.sh` apenas escreve uma mensagem no terminal. A primeira linha indica qual programa deve ser executado ao lê-lo (é preciso executar com `sh` pois não há `bash` na imagem do `alpine`), e a outra linha é o comando que escreve por si. No `Dockerfile`, a primeira linha determina em qual Imagem já existente a nova será criada. Neste caso, será baseado na imagem do `alpine`, pois a imagem `scratch` não possui nem `sh` para executar o script. O segundo comando copia o arquivo de um local tendo referência o arquivo `Dockerfile` para dentro da imagem com destino a raiz (`/`) da imagem. O terceiro comando garante que o `/NossoScript.sh` é executável e a última linha é o comando executado ao executar a imagem. Para construir a imagem, digite o seguinte comando no terminal:

```
$ docker build .
Sending build context to Docker daemon  3.072kB
Step 1/4 : FROM alpine
latest: Pulling from library/alpine
050382585609: Pull complete
Digest: sha256:6a92cd1fcdc8d8cdec60f33dda4db2cb1fcdcacf3410a8e05b3741f44a9b5998
Status: Downloaded newer image for alpine:latest
 ---> b7b28af77ffe
Step 2/4 : COPY NossoScript.sh /
 ---> 4dcaaff92cf5
Step 3/4 : RUN chmod a+x /NossoScript.sh
 ---> Running in da4badb79b54
Removing intermediate container da4badb79b54
 ---> e313578e0ce8
Step 4/4 : CMD ["/NossoScript.sh"]
 ---> Running in f11c2de74409
Removing intermediate container f11c2de74409
 ---> 2030515cbe4c
Successfully built 2030515cbe4c
```

Etapa por etapa o Docker vai construindo a imagem. Primeiramente o Docker verifica se há uma imagem do Alpine na máquina e, caso não haja, a imagem mais recente irá ser baixada. Cada etapa é uma camada na construção da imagem. Após baixar a imagem do Linux Alpine, o Docker marca esta camada com uma hash. Neste caso, a camada com o Alpine foi marcada com a hash `b7b28af77ffe`. As próximas etapas são copiar o script para dentro da imagem e torná-lo executável. Cada etapa marcada com sua respectiva hash. Por fim, é definido o comando que será executado ao iniciar o contêiner e construido a imagem `2030515cbe4c`. Para executá-la, utilize o seguinte comando:

```
$ docker run 2030515cbe4c
Olá do nosso script
```

Identificação de Imagem
-----------------------

Para facilitar a identificação das imagens, existe um parâmetro para colocar um rótulo na imagem construida. É possível rodar o contêiner tanto com a hash da imagem quanto sua _tag_.

```
$ docker build -t 0dockerimage .
Sending build context to Docker daemon  3.072kB
Step 1/4 : FROM alpine
 ---> b7b28af77ffe
Step 2/4 : COPY NossoScript.sh /
 ---> Using cache
 ---> 4dcaaff92cf5
Step 3/4 : RUN chmod a+x /NossoScript.sh
 ---> Using cache
 ---> e313578e0ce8
Step 4/4 : CMD ["/NossoScript.sh"]
 ---> Using cache
 ---> 2030515cbe4c
Successfully built 2030515cbe4c
Successfully tagged 0dockerimage:latest
$ docker run 0dockerimage
Olá do nosso script
```

É comum identificarmos a imagem através de versões, além de seu nome. Isso também pode ser feito com o parâmetro da _tag_ passado no _build_, ou no comando `image tag`
 
```
$ docker image tag 0dockerimage 0dockerimage:v0.0.1
$ docker run 0dockerimage:v0.0.1
Olá do nosso script
```
