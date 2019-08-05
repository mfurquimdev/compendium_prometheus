Docker
======

CMD vs Entrypoint
-----------------

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


```
$ docker build -t 2entrypoint .
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
Successfully tagged 2entrypoint:latest
```

```
$ docker run 2entrypoint:latest
Olá do nosso texto
```

