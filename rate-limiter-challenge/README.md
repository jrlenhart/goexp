# Rate Limiter Challenge

O objetivo deste repositório é resolver o desafio descrito no arquivo [CHALLENGE.md](./CHALLENGE.md)

## Variáveis de ambiente

É possível ajustar as seguintes variáveis de configuração para personalizar o comportamento do rate limiter:

- <b>MAX_REQUESTS_IP:</b> Define o número máximo de requisições permitidas por segundo para um IP. `Default: 5`
- <b>MAX_REQUESTS_TOKEN:</b> Define o número máximo de requisições permitidas por segundo para um token. `Default: 10`
- <b>BLOCK_TIME_IP:</b> Define o tempo em segundos que um IP será bloqueado após exceder o limite de requisições. `Default: 300 (5 minutos)`
- <b>BLOCK_TIME_TOKEN:</b> Define o tempo em segundos que um token será bloqueado após exceder o limite de requisições. `Default: 600 (10 minutos)`

Essas variáveis estão configuradas ao inicializar o servidor e podem ser passadas como argumentos de configuração no arquivo `.env`.

## Executar o projeto

Siga os passos abaixo para executar o projeto localmente:

### 1. Inicie o Docker com o Redis

Execute o comando abaixo para iniciar o Redis utilizando o Docker:

```bash
docker-compose up -d
```

### 2. Execute o servidor

Em seguida, execute o servidor com o comando:

```bash
go run ./cmd/main.go
```

### 3. Realize requisições HTTP

Agora, faça requisições HTTP para testar o rate limiter. Aqui estão dois exemplos utilizando o terminal:

- IP (10 requisições):

```bash
for i in {1..10}; do curl -i localhost:8080; done
```

- Token (12 requisições):

```bash
for i in {1..12}; do curl -i -H "API_KEY: abc1234" localhost:8080; done
```

## Testes

Para executar os testes, utilize o comando abaixo a partir da raiz do projeto:

```bash
go test ./internal/limiter -v
```
