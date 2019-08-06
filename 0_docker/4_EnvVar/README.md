Environment Variables
=====================

Uma forma comum de comunicar-se com um contêiner é utilizando variáveis de ambiente. Esta sessão mostrará como isso pode ser feito.


Declaração de Variáveis
-----------------------

Considere o conteúdo dos seguintes arquivos `printvar.sh` e `Dockerfile`:


**printvar.sh**
```
#!/bin/sh

echo "Olá do script. Variável OLA_MUNDO = $OLA_MUNDO"
```


**Dockerfile**
```
FROM alpine

ENV OLA_MUNDO "Olá do valor padrão"

ADD printvar.sh /
RUN chmod a+x /printvar.sh

ENTRYPOINT ["/printvar.sh"]
```

O arquivo `printvar.sh` escreve uma mensagem e o valor da variável. Neste ponto, o único comando não familiar deve ser o `ENV`. Esta palavra-chave é utilizada para declarar uma variável de ambiente e seu valor. No `Dockerfile` não é preciso declarar a variável de ambiente como foi feito, mas são boas práticas para que melhor entendam o uso da imagem, facilitando assim a manutenção do código. Construa a imagem e rode o contêiner com os comandos abaixo:


```
$ docker build -t 4envvar .
Sending build context to Docker daemon  5.632kB
Step 1/5 : FROM alpine
 ---> b7b28af77ffe
Step 2/5 : ENV OLA_MUNDO "Olá do valor padrão"
 ---> Running in cd3f8376af71
Removing intermediate container cd3f8376af71
 ---> b7ba94ac4d54
Step 3/5 : ADD printvar.sh /
 ---> dccbc83dfe1f
Step 4/5 : RUN chmod a+x /printvar.sh
 ---> Running in ed3fe4f58e22
Removing intermediate container ed3fe4f58e22
 ---> 915ef6a56c54
Step 5/5 : ENTRYPOINT ["/printvar.sh"]
 ---> Running in 62870a678771
Removing intermediate container 62870a678771
 ---> feb7c8c70f09
Successfully built feb7c8c70f09
Successfully tagged 4envvar:latest
$ docker run 4envvar
Olá do script. Variável OLA_MUNDO = Olá do valor padrão
```

Observe que o valor da variável não foi modificado. Para modificá-la, utilize o parâmetro `--env` para o `docker run`:


```
$ docker run --env OLA_MUNDO="Modificando Variável de Ambiente" 4envvar
Olá do script. Variável OLA_MUNDO = Modificando Variável de Ambiente
```
