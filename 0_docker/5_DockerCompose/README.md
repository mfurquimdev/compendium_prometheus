Docker Compose
==============

O `docker-compose` é uma ferramenta que torna mais fácil a execução dos contêiners Docker. Considere o conteúdo dos arquivos `printvar.sh`, `Dockerfile`, e `docker-compose.yml`:


**printvar.sh**
```
#!/bin/sh

echo "Olá do script. Variável OLA_MUNDO=[$OLA_MUNDO]"
```


**Dockerfile**
```
FROM alpine

ENV OLA_MUNDO "Olá do valor padrão"

ADD printvar.sh /
RUN chmod a+x /printvar.sh

ENTRYPOINT ["/printvar.sh"]
```


**docker-compose.yml**
```
version: '3.5'

services:

  envvar:
    build: .
    image: 5dockercompose
    environment:
      - OLA_MUNDO="Alterando pelo docker-compose"
```

Os arquivos `printvar.sh` e `Dockerfile` são os mesmos do exercício passado (`4_EnvVar`). O arquivo novo desta sessão é o `docker-compose.yml` que modifica a forma com que executamos o contêiner. A primeira linha dita a versão do documento<sup>[compose-file-version](https://docs.docker.com/compose/compose-file/compose-versioning/#compatibility-matrix)</sup>, que deve ser compatível com o docker engine da sua máquina. A próxima palavra-chave é _services_ que inicia as definições dos serviços. Cada serviço possui um nome e características únicas, como imagem e variável de ambiente diferentes. Para construir a imagem e subir o contêiner, execute o seguinte comando:

```
$ docker-compose up --build
Building envvar
Step 1/4 : FROM alpine
 ---> b7b28af77ffe
Step 2/4 : ADD printvar.sh /
 ---> a5cdb7774a9a
Step 3/4 : RUN chmod a+x /printvar.sh
 ---> Running in a3034f5e1cc5
Removing intermediate container a3034f5e1cc5
 ---> 5a7eea89df24
Step 4/4 : ENTRYPOINT ["/printvar.sh"]
 ---> Running in 2ca512a7fc7d
Removing intermediate container 2ca512a7fc7d
 ---> 192c9c210e6c

Successfully built 192c9c210e6c
Successfully tagged 5dockercompose:v0.0.1
Recreating 5_dockercompose_envvar_1 ... done
Attaching to 5_dockercompose_envvar_1
envvar_1  | Olá do script. Variável OLA_MUNDO=[Alterando pelo docker-compose]
5_dockercompose_envvar_1 exited with code 0
```

Como a flag `--build` foi passada, o `docker-compose` constroi a imagem primeiro. Após a construção da imagem, o serviço é executado e termina com código 0 (indicando que está tudo certo). O `docker-compose` separa os contêiner com o nome do serviço e o número dele. Neste caso, `envvar` é o nome do serviço e existe apenas um contêiner rodando, por isso o nome `envvar_1`. A saída padrão vem precedido do nome do contêiner, como observado na penúltima linha.

Múltiplos Serviços
------------------

É possível ter mais de um serviço no mesmo `docker-compose.yml`. Considere um `docker-compose.yml` com o seguinte conteúdo:

**docker-compose-multiplos.yml**
```
version: '3.5'

services:

  envvar1:
    build: .
    image: 5dockercompose:v0.0.1
    environment:
      - OLA_MUNDO=Olá do primeiro serviço

  envvar2:
    build: .
    image: 5dockercompose:v0.0.1
    environment:
      - OLA_MUNDO=Olá do segundo serviço
```

Ao executar, temos o seguinte resultado:

```
$ docker-compose -f docker-compose-multiplos.yml up
Creating 5_dockercompose_envvar2_1 ... done
Creating 5_dockercompose_envvar1_1 ... done
Attaching to 5_dockercompose_envvar2_1, 5_dockercompose_envvar1_1
envvar2_1  | Olá do script. Variável OLA_MUNDO=[Olá do segundo serviço]
envvar1_1  | Olá do script. Variável OLA_MUNDO=[Olá do primeiro serviço]
5_dockercompose_envvar2_1 exited with code 0
5_dockercompose_envvar1_1 exited with code 0
```

A flag `-f` é para indicar um outro arquivo que não seja o nome padrão `docker-compose.yml`. Não é necessário construir a imagem pois foi criada no comando anterior. Os dois serviços escrevem o conteúdo de suas respectivcas variáveis de ambiente e terminam com código ok. É possível também executar mais de um contêiner de um serviço. Executando o comando para replicar/escalar, temos o seguinte resultado:

```
$ docker-compose -f docker-compose-multiplos.yml up --scale envvar1=3 --scale envvar2=5
Starting 5_dockercompose_envvar1_1 ... done
Starting 5_dockercompose_envvar2_1 ... done
Creating 5_dockercompose_envvar2_2 ... done
Creating 5_dockercompose_envvar2_3 ... done
Creating 5_dockercompose_envvar2_4 ... done
Creating 5_dockercompose_envvar2_5 ... done
Creating 5_dockercompose_envvar1_2 ... done
Creating 5_dockercompose_envvar1_3 ... done
Attaching to 5_dockercompose_envvar1_1, 5_dockercompose_envvar1_2, 5_dockercompose_envvar1_3, 5_dockercompose_envvar2_1, 5_dockercompose_envvar2_5, 5_dockercompose_envvar2_3, 5_dockercompose_envvar2_2, 5_dockercompose_envvar2_4
envvar1_1  | Olá do script. Variável OLA_MUNDO=[Olá do primeiro serviço]
envvar1_2  | Olá do script. Variável OLA_MUNDO=[Olá do primeiro serviço]
envvar1_3  | Olá do script. Variável OLA_MUNDO=[Olá do primeiro serviço]
envvar2_5  | Olá do script. Variável OLA_MUNDO=[Olá do segundo serviço]
envvar2_1  | Olá do script. Variável OLA_MUNDO=[Olá do segundo serviço]
envvar2_3  | Olá do script. Variável OLA_MUNDO=[Olá do segundo serviço]
envvar2_2  | Olá do script. Variável OLA_MUNDO=[Olá do segundo serviço]
envvar2_4  | Olá do script. Variável OLA_MUNDO=[Olá do segundo serviço]
5_dockercompose_envvar1_1 exited with code 0
5_dockercompose_envvar2_1 exited with code 0
5_dockercompose_envvar2_5 exited with code 0
5_dockercompose_envvar1_2 exited with code 0
5_dockercompose_envvar2_3 exited with code 0
5_dockercompose_envvar1_3 exited with code 0
5_dockercompose_envvar2_2 exited with code 0
5_dockercompose_envvar2_4 exited with code 0
```
