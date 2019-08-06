Docker
======

**Docker** é uma tecnologia que permite criar imagens de contêiner para rodar em diversos ambientes, inclusive na nuvem.

**Imagem Docker** é uma unidade de software que contém o código e todas as suas dependências para ser executado.

**Contêiner** é uma imagem docker em execução.

Os contêiners Docker que rodam no Docker Engine são: o padrão na indústria, podendo ser executado em qualquer local; leves por compartilhar o núcleo do sistema operacional, aumentando a eficiência dos servidores; seguros por ser isolado, rodando cada imagem em um processo separado no espaço do usuário<sup>[docker-what-container](https://www.docker.com/resources/what-container)</sup>.

Estrutura
---------

A primeira coisa neste compêndio é **executar** uma imagem docker qualquer e entender que é possível subir o contêiner em modo interativo ou não. Depois, **construir** uma imagem docker, entendendo seu processo de construção em camadas e a utilização de **cache** para otimizá-lo. Dois conceitos importantes são o de **entrypoint** e **cmd** e, por fim, a comunicação com o contêiner através de **variáveis de ambiente**.
