# Stress Test Challenge

O objetivo deste repositório é resolver o desafio descrito no arquivo [CHALLENGE.md](./CHALLENGE.md)

## Executar o projeto

Para rodar o projeto localmente, siga os passos abaixo:

### 1. Construir a imagem Docker

Execute o seguinte comando para construir a imagem Docker do projeto:

```bash
docker build -t stress-test .
```

### 2. Rodar o teste de carga

Depois, execute o comando abaixo para iniciar o teste de carga:

```bash
docker run stress-test --url=http://google.com --requests=100 --concurrency=10
```

### Parâmetros disponíveis

- `--url`: URL do serviço a ser testado.
- `--requests`: Número total de requests.
- `--concurrency`: Número de chamadas simultâneas.
