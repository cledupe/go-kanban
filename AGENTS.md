# Go Kanban

## Objetivo

MVP de um Kanban com:

- frontend `Next.js`
- backend `Go`
- persistencia `SQLite`
- execucao local e de testes via `Docker` e `docker-compose`

O fluxo esperado e de estado confirmado pelo servidor: a UI so atualiza o estado persistido depois da resposta bem-sucedida da API.

## Status Atual

Referencia de status: `2026-05-20`.

Tasks do plano:

- `task_01` `completed`: bootstrap do backend Go com Docker e Compose
- `task_02` `completed`: persistencia SQLite, migracoes e repositorios base
- `task_03` `pending`: servicos e API REST para boards, columns e cards
- `task_04` `pending`: templates de board e regras de criacao inicial
- `task_05` `pending`: bootstrap do frontend Next.js com cliente HTTP e Compose
- `task_06` `pending`: UI de lista de boards e detalhe do board
- `task_07` `pending`: CRUD de columns e cards com atualizacao confirmada pelo servidor
- `task_08` `pending`: drag-and-drop com persistencia via API

Resumo pratico:

- o backend ja sobe com `App`, `Config`, roteador HTTP e endpoint `GET /readyz`
- o banco SQLite ja abre e executa migracoes no startup
- ja existem entidades de dominio para `Board`, `Column` e `Card`
- ja existem repositorios SQLite com CRUD base
- ainda nao existe camada de servicos
- ainda nao existe API REST do MVP
- ainda nao existe frontend no repositorio

## Estrutura Atual

Arquivos e areas relevantes:

- `backend/cmd/api/main.go`: entrypoint da API
- `backend/internal/app/`: bootstrap da aplicacao, config e ciclo de vida
- `backend/internal/http/`: roteador e readiness
- `backend/internal/domain/`: entidades de dominio
- `backend/internal/storage/sqlite/`: conexao, migracoes e repositorios
- `backend/integration/`: smoke test de container
- `docker-compose.yml`: orquestracao atual do backend

## Arquitetura Proposta

O projeto deve seguir uma arquitetura separada em duas aplicacoes:

1. `frontend`
   - `Next.js`
   - consome a API REST do backend
   - renderiza boards, columns e cards
   - executa CRUD e drag-and-drop
   - so consolida mudancas na UI depois da confirmacao da API

2. `backend`
   - `Go`
   - expoe endpoints REST para boards, columns e cards
   - valida requests
   - concentra regras de negocio na camada de servicos
   - usa repositorios para persistencia

3. `database`
   - `SQLite`
   - armazena boards, columns, cards e ordenacao
   - fica atras de contratos de repositorio para permitir migracao futura

4. `container runtime`
   - `Docker` e `docker-compose`
   - padronizam execucao local, testes e uso por agentes

## Camadas do Backend

Separacao desejada:

- `handlers/http`
  - parse de request
  - validacao estrutural
  - serializacao de response
  - mapeamento de erros de dominio para status HTTP

- `services`
  - regras de negocio
  - validacoes de dominio
  - agregacao de board detail
  - regras de ordenacao e movimentacao

- `repositories`
  - acesso ao SQLite
  - CRUD e queries ordenadas
  - sem regras de negocio

## Contratos HTTP Esperados

Boards:

- `GET /api/boards`
- `POST /api/boards`
- `GET /api/boards/:id`
- `PATCH /api/boards/:id`
- `DELETE /api/boards/:id`

Columns:

- `POST /api/boards/:boardId/columns`
- `PATCH /api/columns/:id`
- `DELETE /api/columns/:id`

Cards:

- `POST /api/columns/:columnId/cards`
- `PATCH /api/cards/:id`
- `POST /api/cards/:id/move`
- `DELETE /api/cards/:id`

Mapeamento esperado de erro:

- `ErrInvalidInput` -> `400`
- `ErrNotFound` -> `404`
- `ErrConflict` -> `409`

## Estado-Alvo do Board Detail

O endpoint `GET /api/boards/:id` deve retornar dados prontos para renderizacao do frontend:

- board
- columns ordenadas por `position`
- cards ordenados por `position` dentro de cada column

Esse contrato e o principal desbloqueador para o frontend inicial.

## Proximo Passo Recomendado

O caminho critico atual e concluir a `task_03`:

- criar erros de dominio e contratos de servico
- implementar servicos de board, column e card
- expor a API REST completa do MVP
- adicionar testes de servico, validacao e status HTTP

## Regras para Agentes

- preserve a separacao entre `http`, `service` e `storage`
- nao mova regra de negocio para repositorios
- mantenha respostas e erros estaveis para o frontend
- siga TDD sempre que estiver adicionando comportamento novo
- antes de marcar algo como concluido, rode verificacao fresca
