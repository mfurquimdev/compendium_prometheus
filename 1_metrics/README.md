Métrica
=======

Métricas de software são medidas de uma característica contável ou quantificável. Após instrumentar um código para expor tais métricas, o Prometheus faz o scrape e armazena todos estes dados de métrica como uma série temporal. É possível criar suas próprias regras e alertas para facilitar a propagação desses dados para um nível superior da federação ou para um outro serviço, como o Grafana.

Tipos de Métrica
----------------

As bibliotecas cliente do Prometheus oferecem quatro tipos de métricas: o **counter**, que é uma métrica cumulativa e representa algo que pode apenas aumentar, como o número de requisições total; o **gauge**, que representa um valor que pode tanto aumentar quanto diminuir, como o número de processos em um determinado momento; o **historgram**, que armazena os valores em determinados "baldes" e expõe uma soma e um contador; e o **summary** que, assim como o histogram, expõe uma soma, um contador e armazena seus dados em "baldes" com valores em porcentagem, e não absolutos.


