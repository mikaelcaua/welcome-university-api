# Welcome University API

## Objetivo do projeto

Esta API organiza a navegacao da estrutura academica e o ciclo de envio e validacao de provas.

O fluxo principal do sistema e:

1. consultar estados, universidades, cursos e disciplinas;
2. autenticar usuarios;
3. permitir envio de provas em PDF vinculadas a uma disciplina;
4. permitir aprovacao ou rejeicao dessas provas por perfis autorizados;
5. disponibilizar apenas provas aprovadas para consulta publica.

## Stack e arquitetura

- Java 17
- Spring Boot 3.2
- Spring Web
- Spring Security com JWT
- Spring Data JPA
- PostgreSQL ou H2
- Springdoc OpenAPI com Swagger UI
- S3 compativel para armazenamento dos arquivos das provas

## Entidades principais

- `State`: unidade federativa usada como ponto inicial da navegacao.
- `University`: universidade vinculada a um estado.
- `Course`: curso vinculado a uma universidade.
- `Subject`: disciplina vinculada a um curso.
- `AppUser`: usuario autenticavel com papel de acesso.
- `Exam`: prova enviada por usuario, associada a uma disciplina e a um arquivo da prova (PDF ou imagem).

## Perfis e permissoes

- `USER`: pode consultar dados publicos, autenticar-se, enviar provas e consultar o proprio perfil.
- `APPROVER`: possui as permissoes de usuario autenticado e pode listar provas pendentes e revisar provas.
- `ADMIN`: possui acesso administrativo de usuarios e tambem pode revisar provas.
- `DEV`: perfil tecnico com privilegios amplos; nao pode ser atribuido pela API.

## Casos de uso

### 1. Navegacao academica publica

Permite que qualquer cliente navegue pela estrutura:

- listar estados;
- buscar um estado por sigla;
- listar universidades de um estado;
- listar cursos de uma universidade;
- listar disciplinas de um curso;
- listar provas aprovadas de forma geral ou por disciplina.

Casos de uso atendidos:

- montar seletores em cascata no frontend;
- exibir catalogo de provas disponiveis;
- buscar provas aprovadas por disciplina e periodo.

### 2. Cadastro e autenticacao

O usuario pode:

- criar conta com nome, email e senha;
- autenticar-se com email e senha;
- renovar sessao com refresh token.

Resultado esperado:

- retorno de `accessToken`;
- retorno de `refreshToken`;
- retorno dos dados basicos do usuario autenticado.

### 3. Envio de prova

Usuario autenticado envia uma prova com:

- nome;
- ano;
- semestre;
- tipo da prova;
- disciplina;
- arquivo da prova (PDF ou imagem).

Fluxo:

1. a API valida o arquivo;
2. valida o usuario autenticado;
3. valida a disciplina;
4. envia o arquivo para o armazenamento S3;
5. cria a prova com status `PENDING`.

### 4. Revisao de prova

Perfis `APPROVER`, `ADMIN` e `DEV` podem:

- listar provas pendentes;
- aprovar prova;
- rejeitar prova;
- registrar observacao de revisao.

Esse fluxo controla a qualidade do acervo visivel publicamente.

### 5. Administracao de usuarios

Perfis `ADMIN` e `DEV` podem:

- listar usuarios;
- alterar papel de um usuario;
- consultar o proprio perfil autenticado.

## Regras de negocio

### Autenticacao e identidade

- O email e normalizado para lowercase e sem espacos laterais no cadastro e no login.
- Nao e permitido cadastrar dois usuarios com o mesmo email.
- O primeiro usuario criado no sistema recebe papel `ADMIN`.
- Todos os usuarios criados depois do primeiro recebem papel inicial `USER`.
- Login com credenciais invalidas retorna erro de autenticacao.
- Refresh token invalido retorna erro de autenticacao.

### Autorizacao

- Endpoints de documentacao Swagger sao publicos.
- Rotas de consulta da arvore academica e de provas aprovadas sao publicas.
- Criacao de estado, universidade e curso exige perfil `ADMIN` ou `DEV`.
- Upload de prova exige autenticacao.
- Revisao de prova exige perfil `APPROVER`, `ADMIN` ou `DEV`.
- Consulta de usuarios exige `ADMIN` ou `DEV`.
- Atualizacao de papel de usuario exige `ADMIN` ou `DEV`.
- O papel `DEV` nao pode ser concedido pela API.

### Regras sobre provas

- Apenas arquivos PDF e imagens sao aceitos no upload.
- Arquivo ausente ou vazio invalida o envio.
- O campo `semester` deve ser `1` ou `2`.
- O campo `examYear` deve ser no minimo `2000`.
- A disciplina informada no upload precisa existir.
- Toda prova enviada entra com status inicial `PENDING`.
- Usuarios com papel `USER` podem ter no maximo 5 provas pendentes ao mesmo tempo.
- Na revisao, nao e permitido definir o status como `PENDING`.
- Apenas provas que ainda estao `PENDING` podem ser revisadas.
- Ao revisar uma prova, a API registra quem revisou, quando revisou e uma observacao opcional.

### Regras de consulta

- A listagem publica retorna apenas provas com status `APPROVED`.
- O filtro `period` so pode ser usado junto com `subjectId`.
- O formato de `period` deve ser `AAAA.S`.
- O semestre dentro de `period` deve ser `1` ou `2`.

## Endpoints por fluxo

### Publicos

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `GET /states`
- `GET /states/{code}`
- `GET /states/{stateId}/universities`
- `GET /universities/{universityId}/courses`
- `GET /courses/{courseId}/subjects`
- `GET /subjects/{subjectId}/exams`
- `GET /exams`

### Autenticados

- `GET /users/me`
- `POST /exams`

### Revisao

- `GET /exams/pending`
- `PATCH /exams/{examId}/status`

### Administracao

- `GET /users`
- `PATCH /users/{id}/role`
- `POST /states`
- `POST /states/{stateId}/universities`
- `POST /universities/{universityId}/courses`

## Payloads e DTOs por rota

### `POST /auth/register`

Autenticacao: publica

Request body:

```json
{
  "name": "Maria Silva",
  "email": "maria@exemplo.com",
  "password": "senha1234"
}
```

DTO de entrada: `RegisterRequest`

- `name: string` obrigatorio
- `email: string` obrigatorio, formato de email
- `password: string` obrigatorio, minimo 8 e maximo 100 caracteres

Response `200/201`:

```json
{
  "accessToken": "jwt-access-token",
  "refreshToken": "jwt-refresh-token",
  "tokenType": "Bearer",
  "expiresInSeconds": 900,
  "user": {
    "id": 1,
    "name": "Maria Silva",
    "email": "maria@exemplo.com",
    "role": "ADMIN",
    "createdAt": "2026-03-04T12:00:00Z"
  }
}
```

DTO de saida: `AuthResponse`

- `accessToken: string`
- `refreshToken: string`
- `tokenType: string`
- `expiresInSeconds: long`
- `user: UserResponse`

### `POST /auth/login`

Autenticacao: publica

Request body:

```json
{
  "email": "maria@exemplo.com",
  "password": "senha1234"
}
```

DTO de entrada: `LoginRequest`

- `email: string` obrigatorio, formato de email
- `password: string` obrigatorio

Response `200`:

Mesmo payload de `AuthResponse` usado em `POST /auth/register`.

### `POST /auth/refresh`

Autenticacao: publica

Request body:

```json
{
  "refreshToken": "jwt-refresh-token"
}
```

DTO de entrada: `RefreshTokenRequest`

- `refreshToken: string` obrigatorio

Response `200`:

Mesmo payload de `AuthResponse` usado em `POST /auth/register`.

### `GET /states`

Autenticacao: publica

Request body: nao possui

Response `200`:

```json
[
  {
    "id": 1,
    "code": "CE",
    "name": "Ceara"
  }
]
```

Retorno: `List<State>`

- `id: long`
- `code: string`
- `name: string`

### `GET /states/{code}`

Autenticacao: publica

Path params:

- `code: string`

Request body: nao possui

Response `200`:

```json
{
  "id": 1,
  "code": "CE",
  "name": "Ceara"
}
```

Retorno: `State`

### `POST /states`

Autenticacao: Bearer token com perfil `ADMIN` ou `DEV`

Request body:

```json
{
  "code": "MA",
  "name": "Maranhao"
}
```

DTO de entrada: `CreateStateRequest`

- `code: string` obrigatorio, 2 caracteres
- `name: string` obrigatorio

Response `201`:

Mesmo payload de `State`.

### `GET /states/{stateId}/universities`

Autenticacao: publica

Path params:

- `stateId: long`

Request body: nao possui

Response `200`:

```json
[
  {
    "id": 10,
    "name": "Universidade Exemplo",
    "state": {
      "id": 1,
      "code": "CE",
      "name": "Ceara"
    }
  }
]
```

Retorno: `List<University>`

- `id: long`
- `name: string`
- `state: State`

### `POST /states/{stateId}/universities`

Autenticacao: Bearer token com perfil `ADMIN` ou `DEV`

Path params:

- `stateId: long`

Request body:

```json
{
  "name": "UFMA"
}
```

DTO de entrada: `CreateUniversityRequest`

- `name: string` obrigatorio

Response `201`:

Mesmo payload de `University`.

### `GET /universities/{universityId}/courses`

Autenticacao: publica

Path params:

- `universityId: long`

Request body: nao possui

Response `200`:

```json
[
  {
    "id": 20,
    "name": "Ciencia da Computacao",
    "university": {
      "id": 10,
      "name": "Universidade Exemplo"
    }
  }
]
```

Retorno: `List<Course>`

- `id: long`
- `name: string`
- `university: University`

### `POST /universities/{universityId}/courses`

Autenticacao: Bearer token com perfil `ADMIN` ou `DEV`

Path params:

- `universityId: long`

Request body:

```json
{
  "name": "Ciencia da Computacao"
}
```

DTO de entrada: `CreateCourseRequest`

- `name: string` obrigatorio

Response `201`:

Mesmo payload de `Course`.

### `GET /courses/{courseId}/subjects`

Autenticacao: publica

Path params:

- `courseId: long`

Request body: nao possui

Response `200`:

```json
[
  {
    "id": 30,
    "name": "Algoritmos",
    "course": {
      "id": 20,
      "name": "Ciencia da Computacao"
    }
  }
]
```

Retorno: `List<Subject>`

- `id: long`
- `name: string`
- `course: Course`

### `GET /subjects/{subjectId}/exams`

Autenticacao: publica

Path params:

- `subjectId: long`

Query params:

- `period: string` opcional, formato `AAAA.S`

Request body: nao possui

Response `200`:

```json
[
  {
    "id": 100,
    "name": "Algoritmos - 2025.1 - PROVA1",
    "examYear": 2025,
    "semester": 1,
    "type": "PROVA1",
    "pdfUrl": "https://bucket.exemplo/provas/arquivo.pdf",
    "status": "APPROVED",
    "subjectId": 30,
    "subjectName": "Algoritmos",
    "uploadedBy": {
      "id": 2,
      "name": "Maria Silva",
      "email": "maria@exemplo.com",
      "role": "USER"
    },
    "reviewedBy": {
      "id": 3,
      "name": "Joao Revisor",
      "email": "joao@exemplo.com",
      "role": "APPROVER"
    },
    "reviewNote": "Arquivo legivel e disciplina correta.",
    "createdAt": "2026-03-04T12:00:00Z",
    "reviewedAt": "2026-03-04T13:00:00Z"
  }
]
```

DTO de saida: `List<ExamResponse>`

### `GET /exams`

Autenticacao: publica

Query params:

- `subjectId: long` opcional
- `period: string` opcional, mas so pode ser usado junto com `subjectId`

Request body: nao possui

Response `200`:

Mesmo payload de `List<ExamResponse>` usado em `GET /subjects/{subjectId}/exams`.

### `GET /exams/pending`

Autenticacao: Bearer token com perfil `APPROVER`, `ADMIN` ou `DEV`

Request body: nao possui

Response `200`:

Mesmo payload de `List<ExamResponse>`, mas retornando provas com status `PENDING`.

### `POST /exams`

Autenticacao: Bearer token

Content-Type:

- `multipart/form-data`

Payload de entrada:

- `examYear: integer` obrigatorio, minimo `2000`
- `semester: integer` obrigatorio, valores `1` ou `2`
- `type: ExamType` obrigatorio
- `subjectId: long` obrigatorio
- `file: PDF ou imagem` obrigatorio

Exemplo de envio multipart:

```text
examYear=2025
semester=1
type=PROVA1
subjectId=30
file=<arquivo.pdf ou arquivo.png>
```

DTO de entrada: `ExamUploadRequest`

Response `201`:

Mesmo payload de `ExamResponse`.

Observacao:

- o valor exato aceito em `type` depende do enum `ExamType` implementado no projeto.
- o campo `name` nao e enviado pela API cliente; ele e gerado pelo backend a partir de disciplina, periodo e tipo da prova.

### `PATCH /exams/{examId}/status`

Autenticacao: Bearer token com perfil `APPROVER`, `ADMIN` ou `DEV`

Path params:

- `examId: long`

Request body:

```json
{
  "status": "APPROVED",
  "reviewNote": "Arquivo valido."
}
```

DTO de entrada: `ExamReviewRequest`

- `status: ExamStatus` obrigatorio
- `reviewNote: string` opcional

Valores relevantes de `ExamStatus`:

- `PENDING`
- `APPROVED`
- `REJECTED`

Response `200`:

Mesmo payload de `ExamResponse`.

### `GET /users/me`

Autenticacao: Bearer token

Request body: nao possui

Response `200`:

```json
{
  "id": 2,
  "name": "Maria Silva",
  "email": "maria@exemplo.com",
  "role": "USER",
  "createdAt": "2026-03-04T12:00:00Z"
}
```

DTO de saida: `UserResponse`

- `id: long`
- `name: string`
- `email: string`
- `role: Role`
- `createdAt: Instant`

### `GET /users`

Autenticacao: Bearer token com perfil `ADMIN` ou `DEV`

Request body: nao possui

Response `200`:

```json
[
  {
    "id": 2,
    "name": "Maria Silva",
    "email": "maria@exemplo.com",
    "role": "USER",
    "createdAt": "2026-03-04T12:00:00Z"
  }
]
```

DTO de saida: `List<UserResponse>`

### `PATCH /users/{id}/role`

Autenticacao: Bearer token com perfil `ADMIN` ou `DEV`

Path params:

- `id: long`

Request body:

```json
{
  "role": "APPROVER"
}
```

DTO de entrada: `UpdateUserRoleRequest`

- `role: Role` obrigatorio

Valores de `Role`:

- `USER`
- `APPROVER`
- `ADMIN`
- `DEV`

Response `200`:

Mesmo payload de `UserResponse`.

Observacao:

- embora `DEV` exista no enum, esse papel nao pode ser atribuido por essa rota.

## Fluxo resumido de uso

### Fluxo publico

1. cliente consulta estados;
2. escolhe universidade, curso e disciplina;
3. consulta provas aprovadas.

### Fluxo de contribuicao

1. usuario cria conta ou faz login;
2. recebe token JWT;
3. envia prova em PDF;
4. prova fica pendente ate revisao.

### Fluxo de moderacao

1. aprovador lista provas pendentes;
2. analisa a submissao;
3. aprova ou rejeita;
4. somente provas aprovadas aparecem na consulta publica.

## Documentacao e testes manuais

- Swagger UI: `/swagger-ui/index.html`
- OpenAPI JSON: `/v3/api-docs`

Para testar endpoints protegidos no Swagger:

1. obtenha um `accessToken` em `/auth/login` ou `/auth/register`;
2. clique em `Authorize`;
3. informe o token Bearer;
4. execute as rotas protegidas com o contexto autenticado.
