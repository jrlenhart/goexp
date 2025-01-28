<b>Objetivo:</b> Criar um sistema CLI em Go para realizar testes de carga em um serviço web. O usuário deverá fornecer a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas.

O sistema deverá gerar um relatório com informações específicas após a execução dos testes.

<b>Entrada de Parâmetros via CLI:</b>

<b>--url:</b> URL do serviço a ser testado.<br>
<b>--requests:</b> Número total de requests.<br>
<b>--concurrency:</b> Número de chamadas simultâneas.

<b>Execução do Teste:</b>

- Realizar requests HTTP para a URL especificada.
- Distribuir os requests de acordo com o nível de concorrência definido.
- Garantir que o número total de requests seja cumprido.

<b>Geração de Relatório:</b>

- Apresentar um relatório ao final dos testes contendo:
    - Tempo total gasto na execução.
    - Quantidade total de requests realizados.
    - Quantidade de requests com status HTTP 200.
    - Distribuição de outros códigos de status HTTP (como 404, 500, etc.).

1. <b>Execução da aplicação:</b>

- Poderemos utilizar essa aplicação fazendo uma chamada via docker. Ex:
    - docker run <sua imagem docker\> --url=http://google.com --requests=1000 --concurrency=10
