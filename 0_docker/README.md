Docker
======

**Docker** é uma tecnologia que permite criar imagens de contêiner para rodar em diversos ambientes, inclusive na nuvem.

**Imagem Docker** é uma unidade de software que contém o código e todas as suas dependências para ser executado.

**Contêiner** é uma imagem docker em execução.

Os contêiners Docker que rodam no Docker Engine são: o padrão na indústria, podendo ser executado em qualquer local; leves por compartilhar o núcleo do sistema operacional, aumentando a eficiência dos servidores; seguros por ser isolado, rodando cada imagem em um processo separado no espaço do usuário<sup>[docker-what-container](https://www.docker.com/resources/what-container)</sup>.

Estrutura
---------

A primeira coisa neste compêndio é **executar** [0\_HelloWorld] uma imagem docker qualquer e entender que é possível subir o contêiner em modo interativo ou não. Depois, **construir** [1\_DockerImage] uma imagem docker, entendendo seu processo de construção em camadas e a utilização de **cache** [2\_Cache] para otimizá-lo. Dois conceitos importantes são o de **entrypoint** e **cmd** [3\_Entrypoint] e a comunicação com o contêiner através de **variáveis de ambiente** [4\_EnvVar]. Por fim, utilizar o **docker-compose** [5\_DockerCompose] para facilitar rodar os contêiners.


Troubleshoot
------------

Se estiver tendo problemas com o proxy, adicione `"bip" : "192.168.233.1/24"` no `~/.docker/daemon.json`:

**~/.docker/daemon.json**
```
{
  "bip" : "192.168.233.1/24",
  "debug" : true,
  "experimental" : false
}
```

TODO
----

* [ ] Explicar `VOLUMES [""]`
* [ ] Explicar `EXPOSE`
* [ ] Explicar `FROM <image> AS <alias>`
