# Analise do Repositorio

## Stack atual

- Linguagem: Java 17
- Framework: Spring Boot 3.2.0
- Persistencia: Spring Data JPA
- Bancos declarados: PostgreSQL e H2
- Documentacao da API: Springdoc OpenAPI

## Estrutura observada

- A aplicacao e pequena e centrada em entidades JPA simples: `State`, `University`, `Course`, `Subject` e `Exam`.
- Os endpoints atuais sao somente de leitura.
- Nao havia pasta `src/main/resources`, entao as configuracoes de banco e ambiente nao estavam externalizadas.
- Nao havia `Dockerfile`, `docker-compose.yml`, workflow de CI/CD nem documentacao operacional.
- O repositorio possui artefatos compilados versionados em `target/`, o que tende a gerar ruido em commit e diff.

## Problemas e riscos encontrados

### 1. Seed rodava em todo startup

O `DataLoader` inseria dados sempre que a aplicacao subia. Em ambiente com volume persistente isso criaria duplicidade ou erro de consistencia com qualquer reinicio de container.

Tratamento aplicado:

- o seed agora depende de `APP_SEED_ENABLED`;
- se a base ja tiver registros, o seed e ignorado.

### 2. Ausencia de configuracao por ambiente

O projeto declarava dependencias de `H2` e `PostgreSQL`, mas nao havia arquivo de configuracao em `resources`.

Tratamento aplicado:

- criacao de [application.yml](/Users/mikael/Documents/welcome-university-api/src/main/resources/application.yml);
- defaults locais com H2;
- parametros de producao lidos por variaveis de ambiente no container.

### 3. Sem readiness/liveness para orquestracao

Nao existia endpoint operacional para healthcheck.

Tratamento aplicado:

- adicionado `spring-boot-starter-actuator`;
- expostos apenas `health` e `info`;
- `Dockerfile` e `docker-compose.yml` usam readiness probe.

### 4. Sem pipeline de entrega

Nao havia automacao para build e distribuicao da imagem.

Tratamento aplicado:

- criacao de [container-publish.yml](/Users/mikael/Documents/welcome-university-api/.github/workflows/container-publish.yml) para disparar no `push` da `main`.

### 5. Artefatos compilados versionados

Arquivos dentro de `target/` estao rastreados pelo Git. Isso nao e uma boa pratica para um projeto Java e polui o historico com binarios gerados.

Tratamento parcial aplicado:

- adicionado [.gitignore](/Users/mikael/Documents/welcome-university-api/.gitignore) para evitar novos arquivos temporarios e builds locais acidentais;
- a remocao dos artefatos ja versionados deve ser feita de forma consciente em um commit de limpeza separado.

## Lacunas que continuam existindo

- Nao ha autenticacao/autorizacao. Hoje a API continua publica dentro do stack.
- Nao ha migrations versionadas. O projeto ainda depende de `ddl-auto=update`.
- Nao ha testes automatizados no repositorio.
- Nao ha backup automatizado do PostgreSQL no proprio projeto.
- Nao ha proxy reverso/TLS no stack desta aplicacao. Isso deve ser tratado no homelab com Traefik, Caddy, Nginx Proxy Manager ou equivalente.

## Prioridades recomendadas depois desta etapa

1. Introduzir migrations com Flyway antes de crescer o schema.
2. Adicionar testes de integracao para os endpoints principais.
3. Definir estrategia de autenticacao se a API nao for estritamente interna.
4. Acoplar backup do volume/banco na rotina do homelab.
