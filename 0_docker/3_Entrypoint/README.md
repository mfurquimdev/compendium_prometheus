CMD vs ENTRYPOINT
=================

O Docker define duas palavras chaves que permitem maior versatilidade na construção de imagens. `CMD` e `ENTRYPOINT` é o conjunto que define o que será executado ao rodar o contêiner da imagem. Considere o conteúdo dos seguintes arquivos `NossoTexto.txt` e `Dockerfile`:

**NossoTexto.txt**
```
Olá do nosso texto
```

**Dockerfile**
```
FROM alpine

COPY NossoTexto.txt /

ENTRYPOINT ["/bin/cat"]
CMD ["/NossoTexto.txt"]
```

O arquivo `NossoTexto.txt` nada mais é do que texto, nada de especial. O `Dockerfile`, por outro lado, é onde se encontra a complexidade. Desta vez, o arquivo não vai ser executado por um `shell`, mas sim terá o seu conteúdo conCATenado na saída padrão (terminal). O programa definido em `ENTRYPOINT` é o programa executado ao rodar o contêiner. Tudo o que estiver em  `CMD` é o argumento passado para o `ENTRYPOINT`. Por padrão, o `ENTRYPOINT` é um `shell`<sup>[stackoverflow](https://stackoverflow.com/questions/21553353/what-is-the-difference-between-cmd-and-entrypoint-in-a-dockerfile)</sup>. Neste caso, o programa a ser executado é o `cat` e o argumento passado a ele é o `NossoTexto.txt`. Como resultado, ao subir o contêiner, o conteúdo _Olá do nosso texto_ é escrito na tela. Construa e execute a imagem com os comandos abaixo:

```
$ docker build -t 3entrypoint .
Sending build context to Docker daemon  6.144kB
Step 1/4 : FROM alpine
 ---> b7b28af77ffe
Step 2/4 : COPY NossoTexto.txt /
 ---> a71105a3a2b2
Step 3/4 : ENTRYPOINT ["/bin/cat"]
 ---> Running in 21cb40753df6
Removing intermediate container 21cb40753df6
 ---> 71c21625d4a3
Step 4/4 : CMD ["/NossoTexto.txt"]
 ---> Running in 8aaaaa856b55
Removing intermediate container 8aaaaa856b55
 ---> f70742959007
Successfully built f70742959007
Successfully tagged 3entrypoint:latest
$ docker run 3entrypoint:latest
Olá do nosso texto
```

