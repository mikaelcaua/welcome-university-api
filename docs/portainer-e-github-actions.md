# Portainer e GitHub Actions

## Pipeline criada

Arquivo: [container-publish.yml](/Users/mikael/Documents/welcome-university-api/.github/workflows/container-publish.yml)

Ela executa quando houver `push` na branch `main` e tambem por disparo manual.

## O que a workflow faz

1. Faz checkout do repositorio.
2. Configura Java 17.
3. Gera o jar com Maven.
4. Roda scan de seguranca no repositorio com Trivy.
5. Faz login no GHCR.
6. Construi e publica a imagem Docker.
7. Roda scan de seguranca na imagem publicada.
8. Opcionalmente dispara um webhook do Portainer.

## Tags publicadas

- `ghcr.io/<owner>/<repo>:latest`
- `ghcr.io/<owner>/<repo>:sha-<commit>`

## Secrets esperados

### Obrigatorios

Nenhum extra alem do `GITHUB_TOKEN` padrao do Actions.

### Opcionais

- `PORTAINER_WEBHOOK_URL`: se definido, a workflow chama o webhook ao final do publish.

## Como usar no Portainer

### Opcao recomendada

1. Publique a imagem via GitHub Actions.
2. No Portainer, crie uma Stack usando [docker-compose.yml](/Users/mikael/Documents/welcome-university-api/docker-compose.yml).
3. Defina as variaveis diretamente na Stack ou use upload de arquivo de ambiente fora do Git.
4. Ajuste `APP_IMAGE` para `ghcr.io/<owner>/<repo>:latest`.
5. Se o pacote no GHCR estiver privado, configure credencial de registry no Portainer.

### Atualizacao automatica

Se quiser que cada `push` na `main` acione o redeploy:

1. Gere um webhook de update da stack no Portainer.
2. Salve a URL em `PORTAINER_WEBHOOK_URL` nos secrets do GitHub.
3. Mantenha a stack apontando para a tag `latest` ou troque para tag fixa quando quiser controle manual.

## Fluxo final

`push na main` -> `GitHub Actions` -> `build + scan + push para GHCR` -> `webhook opcional do Portainer` -> `stack atualizada no homelab`

## Ponto de atencao

Se voce preferir deploy imutavel, nao use `latest` no Portainer. Nesse caso, troque a stack para uma tag fixa e automatize a substituicao do compose em outro passo de CI/CD.
