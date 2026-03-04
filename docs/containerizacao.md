# Containerizacao

## Arquitetura do stack

O deploy foi estruturado em dois containers:

- `api`: imagem da aplicacao Spring Boot publicada no GHCR.
- `db`: imagem oficial `postgres:16-alpine` com volume persistente.

Essa separacao e a forma correta para homelab e para Portainer. O banco esta no stack em sua propria imagem, sem acoplar processo de banco dentro do mesmo container da API.

## Arquivos criados

- [Dockerfile](/Users/mikael/Documents/welcome-university-api/Dockerfile): build multi-stage com Maven e runtime Java 17 em usuario nao root.
- [docker-compose.yml](/Users/mikael/Documents/welcome-university-api/docker-compose.yml): stack principal para deploy.
- [docker-compose.local.yml](/Users/mikael/Documents/welcome-university-api/docker-compose.local.yml): override para build local da imagem a partir do codigo.
- [.env.example](/Users/mikael/Documents/welcome-university-api/.env.example): variaveis de ambiente esperadas.

## Variaveis de ambiente

O arquivo [.env.example](/Users/mikael/Documents/welcome-university-api/.env.example) deve continuar no repositorio como referencia das variaveis esperadas.

Se voce usar Portainer, pode preencher essas variaveis diretamente na Stack e nao precisa criar `.env` no repositorio.

Se voce rodar localmente fora do Portainer, copie `.env.example` para `.env` e ajuste:

- `APP_IMAGE`: imagem publicada no GHCR.
- `APP_PORT`: porta exposta da API no host.
- `APP_SEED_ENABLED`: `true` apenas no primeiro bootstrap se quiser popular dados de exemplo.
- `JAVA_OPTS`: limites e parametros da JVM.
- `SPRING_JPA_HIBERNATE_DDL_AUTO`: hoje o padrao sugerido e `update`.
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`

## Comportamento da aplicacao em container

- A API usa `PostgreSQL` quando recebe as variaveis do compose.
- O endpoint de health fica em `/actuator/health/readiness`.
- O seed de dados fica desligado por padrao.
- O sistema de arquivos do container da API fica `read_only`, com `tmpfs` em `/tmp`.

## Fluxos de uso

### 1. Deploy padrao no homelab

Use [docker-compose.yml](/Users/mikael/Documents/welcome-university-api/docker-compose.yml) no Portainer ou no host Docker.

No Portainer:

- cole o conteudo do compose na Stack;
- preencha as variaveis de ambiente na propria interface;
- nao suba `.env` real para o Git.

### 2. Build local a partir do codigo

Use os dois arquivos:

```bash
cp .env.example .env
docker compose -f docker-compose.yml -f docker-compose.local.yml up --build -d
```

### 3. Subida usando imagem publicada

```bash
cp .env.example .env
docker compose up -d
```

## Persistencia

- O PostgreSQL usa o volume nomeado `postgres_data`.
- A API nao depende de volume proprio.
- Backups devem ser feitos a partir do container `db` ou do volume do banco.

## Observacoes operacionais

- Nao expose a porta do banco externamente se nao houver necessidade.
- Em Portainer, prefira stack em rede interna e exponha apenas a API para o proxy reverso.
- Em restart da stack, o banco sobe primeiro e a API espera o healthcheck do PostgreSQL.
