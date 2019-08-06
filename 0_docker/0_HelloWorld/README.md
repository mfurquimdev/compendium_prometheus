Hello World
===========

Para aprender a tecnologia Docker e criar sua própria imagem, é melhor começar executando uma imagem _hello-world_ e uma imagem em modo interativo.

Imagem Hello World
------------------

Para executar a imagem de Olá Mundo do docker, digite o seguinte comando no terminal:

```
$ docker run hello-world
```

Caso o Docker não encontre a imagem local, a imagem mais recente do `hello-world` será baixada. Após garantir que a imagem local é a mais recente, o Docker executa a imagem e um contêiner é criado, executando o programa `hello` dentro da imagem (um programa compilado do fonte `hello.c`<sup>[hello-world](https://github.com/docker-library/hello-world)</sup>, que escreve apenas algumas instruções na tela. Após escrever na tela, o script termina e o contêiner é terminado.

Execucão em modo interativo
---------------------------

É possível executar alguns comandos de forma interativa, como um `shell`, por exemplo. Execute o seguinte comando no terminal para executar uma imagem do Linux Alpine (possui apenas 5MB)<sup>[alpine-github](https://github.com/alpinelinux/docker-alpine)</sup>:

```
$ docker run -it alpine sh
```

O contêiner não é terminado pois o programa rodando é o `sh` que, por sua vez, ainda está rodando. Quando este programa terminar (digitando `exit` ou ^D), o contêiner também irá terminar.
