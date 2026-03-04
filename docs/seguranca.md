# Seguranca

## Medidas aplicadas nesta entrega

- Container da API roda com usuario nao root.
- `read_only: true` na API, com `tmpfs` para escrita temporaria.
- `cap_drop: ALL` e `no-new-privileges:true` na API.
- Apenas endpoints `health` e `info` do Actuator ficam expostos.
- Scan de seguranca no repositorio e na imagem via Trivy no GitHub Actions.
- Segredos fora do repositorio por meio de `.env`, com [.gitignore](/Users/mikael/Documents/welcome-university-api/.gitignore) impedindo commit acidental.
- Banco isolado em rede interna do compose, sem publicacao de porta para o host.

## O que voce ainda deve fazer no homelab

1. Trocar `POSTGRES_PASSWORD` por uma senha forte e unica.
2. Colocar a API atras de um proxy reverso com TLS.
3. Restringir acesso ao Portainer e ao host Docker.
4. Fazer backup recorrente do banco e testar restore.
5. Manter imagens atualizadas, especialmente `postgres:16-alpine` e a base Java da API.

## Riscos atuais do backend

### Sem autenticacao

Hoje qualquer cliente que alcance a API consegue consultar os endpoints. Se esse backend sair da rede interna, isso vira problema imediato.

### Sem migrations versionadas

`ddl-auto=update` ajuda a subir rapido, mas nao e uma estrategia robusta de producao de longo prazo.

### Sem suite de testes

Nao existe protecao automatica contra regressao funcional no pipeline.

## Recomendacoes praticas

### Curto prazo

- manter a API acessivel apenas pela rede interna ou por VPN;
- desligar `APP_SEED_ENABLED` depois da carga inicial;
- usar tag fixa da imagem em ambientes mais sensiveis.

### Medio prazo

- adotar Flyway;
- adicionar autenticacao se houver consumo externo;
- instrumentar logs e metricas em uma stack centralizada.
