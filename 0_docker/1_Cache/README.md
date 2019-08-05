Docker
======

Cache
-----

Devido ao processo de criação de imagem otimizado do Docker, as imagens são criadas por camadas, começando pelo topo do arquivo `Dockerfile` até o final. Caso não haja alteração nas primeiras linhas, o Docker apenas reutiliza a camada já existente. Se alterarmos o arquivo adicionando um comando echo, ao construir a imagem nova, as camadas de cima que não sofreram alterações serão reaproveitadas, reduzindo o tempo de construção da imagem.

**Dockerfile**
```
FROM alpine

COPY NossoScript.sh /
RUN chmod a+x /NossoScript.sh

RUN echo "Olá do momento de Build"

CMD ["/NossoScript.sh"]
```

Ao executar o comando de construir imagem com o novo `Dockerfile` e rodar o contêiner, temos os seguintes resultados:

```
$ docker build -t 1cache:v0.0.1 .
Sending build context to Docker daemon  3.072kB
Step 1/5 : FROM alpine
 ---> b7b28af77ffe
Step 2/5 : COPY NossoScript.sh /
 ---> Using cache
 ---> 4dcaaff92cf5
Step 3/5 : RUN chmod a+x /NossoScript.sh
 ---> Using cache
 ---> e313578e0ce8
Step 4/5 : RUN echo "Olá do momento de Build"
 ---> Running in 3e3f6246bb87
Olá do momento de Build
Removing intermediate container 3e3f6246bb87
 ---> e5a2869ec624
Step 5/5 : CMD ["/NossoScript.sh"]
 ---> Running in e67411ca967e
Removing intermediate container e67411ca967e
 ---> 08f9291ac461
Successfully built 08f9291ac461
Successfully tagged 1cache:v0.0.1
$ docker run 1cache:v0.0.1
Olá do nosso script
```


Momento de construção vs Momento de execução
--------------------------------------------

Como pode ser observado na sessão anterior, há uma diferença nos comandos executados em momento de construção da imagem e momento de execução do contêiner. Ao executar o `docker build`, o comando `RUN echo "Olá do momento de Build"` é executado. Este comando é executado apenas neste momento de construção de imagem e, apesar de não modificar efetivamente a imagem, o Docker considera esta etapa como uma camada da imagem. Ao rodar o contêiner desta imagem, este comando não é executado novamente. Da mesma forma, o script `NossoScript.sh` é executado apenas na execução do contêiner, e não na construção da imagem.

