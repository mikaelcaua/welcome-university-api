# Documentacao de Deploy e Operacao

Esta pasta centraliza a analise do backend e a estrutura criada para containerizacao, publicacao da imagem e deploy no homelab.

## Arquivos

- [analise-repositorio.md](/Users/mikael/Documents/welcome-university-api/docs/analise-repositorio.md): leitura tecnica do estado atual do projeto e riscos identificados.
- [containerizacao.md](/Users/mikael/Documents/welcome-university-api/docs/containerizacao.md): arquitetura dos containers, variaveis, build e operacao via Docker Compose.
- [portainer-e-github-actions.md](/Users/mikael/Documents/welcome-university-api/docs/portainer-e-github-actions.md): fluxo de CI/CD com GitHub Actions e publicacao/atualizacao no Portainer.
- [seguranca.md](/Users/mikael/Documents/welcome-university-api/docs/seguranca.md): decisoes de hardening e checklist de seguranca para o homelab.

## Arquivos de infraestrutura adicionados na raiz

- [Dockerfile](/Users/mikael/Documents/welcome-university-api/Dockerfile)
- [docker-compose.yml](/Users/mikael/Documents/welcome-university-api/docker-compose.yml)
- [docker-compose.local.yml](/Users/mikael/Documents/welcome-university-api/docker-compose.local.yml)
- [.env.example](/Users/mikael/Documents/welcome-university-api/.env.example)
- [.github/workflows/container-publish.yml](/Users/mikael/Documents/welcome-university-api/.github/workflows/container-publish.yml)

## Resumo rapido

- A API agora esta preparada para subir em container com Java 17 e healthcheck do Spring Actuator.
- O banco entra no stack como uma imagem dedicada do `postgres:16-alpine`, persistindo dados em volume nomeado.
- O `push` na branch `main` dispara build, scan de seguranca, publicacao no GHCR e, opcionalmente, webhook do Portainer.
- O seed de dados deixou de rodar sempre: ele agora depende da variavel `APP_SEED_ENABLED` e nao duplica registros em reinicios.

