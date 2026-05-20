# TechSpec: Go Kanban

## Executive Summary

Go Kanban MVP will be implemented as two separate applications: a Next.js frontend responsible for the board UI and user interactions, and a Go backend responsible for board, column, and card APIs plus persistence. The backend will expose a small REST API consumed by the frontend. Persistence will use SQLite for the initial MVP to reduce operational overhead, while the internal repository contracts will be designed to allow later migration to PostgreSQL without rewriting the service layer.

To keep local development, automated tests, and agent-driven execution consistent across Codex, Claude, and other ACP environments, the project will ship with containerized workflows for both applications. Docker will provide reproducible runtime images, and Docker Compose will orchestrate frontend and backend services with stable ports, mounted source code, and environment variables. The primary trade-off is operational simplicity versus future scalability. SQLite keeps the MVP easy to run and test, but it constrains concurrent writes and requires discipline in repository boundaries to avoid a costly migration later.

## System Architecture

### Component Overview

- `frontend` (Next.js)
  - Renders board list, board detail view, columns, cards, and CRUD interactions.
  - Owns drag-and-drop interaction and form state.
  - Calls backend REST endpoints over HTTP.
  - Waits for API confirmation before mutating visible state for persisted actions.

- `backend` (Go API)
  - Exposes REST endpoints for boards, columns, and cards.
  - Validates requests and orchestrates domain operations.
  - Applies ordering and movement rules for columns and cards.
  - Persists state through repository interfaces.

- `sqlite database`
  - Stores boards, columns, cards, and ordering metadata.
  - Uses a local file-based database for MVP durability.

- `docker runtime`
  - Standardizes execution across developer machines and AI agent environments.
  - Packages frontend and backend with pinned base images and startup commands.

- `docker compose`
  - Orchestrates local multi-service startup for development and test flows.
  - Provides stable service names, shared network, and bind mounts where appropriate.

### Data Flow

1. User interacts with the Next.js UI.
2. Frontend sends HTTP request to Go backend.
3. Handler validates and maps the request into service calls.
4. Service applies business rules and calls repositories.
5. Repository reads or writes SQLite.
6. Backend returns normalized JSON response.
7. Frontend updates UI only after a successful response.

### External System Interactions

- No third-party services in MVP.
- The only system boundaries are HTTP between frontend and backend, SQL access from backend to SQLite, and container orchestration through Docker Compose.

## Implementation Design

### Core Interfaces

```go
type BoardService interface {
    ListBoards(ctx context.Context) ([]Board, error)
    CreateBoard(ctx context.Context, input CreateBoardInput) (Board, error)
    GetBoard(ctx context.Context, boardID string) (BoardDetail, error)
    DeleteBoard(ctx context.Context, boardID string) error
}
```

```go
type CardRepository interface {
    Create(ctx context.Context, card Card) (Card, error)
    Update(ctx context.Context, card Card) (Card, error)
    Delete(ctx context.Context, cardID string) error
    Move(ctx context.Context, cardID string, targetColumnID string, position int) error
}
```

```go
type Card struct {
    ID          string
    ColumnID    string
    Title       string
    Description string
    Position    int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

Error handling conventions:
- Services return domain errors such as `ErrNotFound`, `ErrInvalidInput`, and `ErrConflict`.
- Handlers map domain errors to HTTP status codes.
- Repositories return wrapped storage errors with operation context.

### Data Models

#### Domain Entities

- `Board`
  - `id`
  - `name`
  - `created_at`
  - `updated_at`

- `Column`
  - `id`
  - `board_id`
  - `name`
  - `position`
  - `created_at`
  - `updated_at`

- `Card`
  - `id`
  - `column_id`
  - `title`
  - `description`
  - `position`
  - `created_at`
  - `updated_at`

#### Storage Schema

`boards`
- `id TEXT PRIMARY KEY`
- `name TEXT NOT NULL`
- `created_at DATETIME NOT NULL`
- `updated_at DATETIME NOT NULL`

`columns`
- `id TEXT PRIMARY KEY`
- `board_id TEXT NOT NULL`
- `name TEXT NOT NULL`
- `position INTEGER NOT NULL`
- `created_at DATETIME NOT NULL`
- `updated_at DATETIME NOT NULL`
- foreign key (`board_id`) references `boards(id)` on delete cascade

`cards`
- `id TEXT PRIMARY KEY`
- `column_id TEXT NOT NULL`
- `title TEXT NOT NULL`
- `description TEXT NOT NULL DEFAULT ''`
- `position INTEGER NOT NULL`
- `created_at DATETIME NOT NULL`
- `updated_at DATETIME NOT NULL`
- foreign key (`column_id`) references `columns(id)` on delete cascade

#### API DTOs

- `CreateBoardRequest`
  - `name string`
  - `template string optional`

- `CreateColumnRequest`
  - `name string`

- `CreateCardRequest`
  - `title string`
  - `description string optional`

- `MoveCardRequest`
  - `target_column_id string`
  - `position int`

### API Endpoints

#### Boards

- `GET /api/boards`
  - Returns all boards for the single-user workspace.
  - `200 OK`

- `POST /api/boards`
  - Creates a board from a blank layout or a template.
  - Request: `{ "name": "...", "template": "basic-kanban|bug-tracker|content-pipeline" }`
  - `201 Created`, `400 Bad Request`

- `GET /api/boards/:id`
  - Returns board details including ordered columns and cards.
  - `200 OK`, `404 Not Found`

- `PATCH /api/boards/:id`
  - Renames a board.
  - `200 OK`, `400 Bad Request`, `404 Not Found`

- `DELETE /api/boards/:id`
  - Deletes a board and its dependent columns and cards.
  - `204 No Content`, `404 Not Found`

#### Columns

- `POST /api/boards/:boardId/columns`
  - Adds a column to a board.
  - `201 Created`, `400 Bad Request`, `404 Not Found`

- `PATCH /api/columns/:id`
  - Renames or reorders a column.
  - `200 OK`, `400 Bad Request`, `404 Not Found`

- `DELETE /api/columns/:id`
  - Deletes a column and its cards.
  - `204 No Content`, `404 Not Found`

#### Cards

- `POST /api/columns/:columnId/cards`
  - Creates a card in a column.
  - `201 Created`, `400 Bad Request`, `404 Not Found`

- `PATCH /api/cards/:id`
  - Updates title and description.
  - `200 OK`, `400 Bad Request`, `404 Not Found`

- `POST /api/cards/:id/move`
  - Moves a card to another column and position.
  - `200 OK`, `400 Bad Request`, `404 Not Found`, `409 Conflict`

- `DELETE /api/cards/:id`
  - Deletes a card.
  - `204 No Content`, `404 Not Found`

### Containerization

#### Dockerfiles

- `frontend/Dockerfile`
  - Multi-stage build for dependency install, app build, and runtime image.
  - Development target supports mounted source and `next dev`.
  - Test target supports CI-style execution in isolated environments.

- `backend/Dockerfile`
  - Multi-stage build for Go dependency resolution, test execution, and final binary image.
  - Development target supports mounted source and hot-reload wrapper if added later.
  - Test target runs unit and integration tests consistently in local and agent environments.

#### Docker Compose Services

- `frontend`
  - Exposes Next.js on a fixed local port.
  - Receives backend base URL through environment variables.
  - Mounts source code in development mode.

- `backend`
  - Exposes Go API on a fixed local port.
  - Mounts SQLite data directory or local volume for persistence.
  - Runs migrations on startup before serving requests.

- `test`
  - Optional compose profile or dedicated service for standardized test execution.
  - Supports commands such as backend tests, frontend tests, and end-to-end checks.

#### Expected Compose Responsibilities

- Provide one command to boot the full MVP stack locally.
- Provide one command to run reproducible test workflows.
- Avoid assumptions about host-installed Go or Node versions.
- Work in agent environments where only Docker is reliably available.

## Integration Points

No external integrations are required for MVP.

## Impact Analysis

| Component | Impact Type | Description and Risk | Required Action |
|-----------|-------------|---------------------|-----------------|
| Frontend Next.js app | new | Entire user-facing application must be created; moderate delivery risk | Build board list, board detail, forms, and drag-and-drop UI |
| Go API service | new | Entire backend service must be created; moderate delivery risk | Implement handlers, services, repositories, and migrations |
| SQLite schema | new | New persistence layer with ordering logic; moderate data integrity risk | Define schema, migrations, and transactional move operations |
| Template module | new | Required to seed default board structures; low risk | Implement static template definitions in backend |
| Dockerfiles | new | Reproducible environment depends on correct container boundaries; low risk | Add frontend and backend Dockerfiles with dev and test targets |
| Docker Compose | new | Misconfigured volumes or ports can slow onboarding; low risk | Add compose file for local run and test orchestration |
| Test harness | new | No existing coverage or fixtures; moderate quality risk | Add unit and integration tests for core flows |

## Testing Approach

### Unit Tests

- Service-layer tests for:
  - board creation from template
  - board rename and delete
  - column create, rename, reorder, delete
  - card create, edit, delete
  - card move across columns
- Repository tests for:
  - position assignment
  - move semantics
  - cascade deletion behavior
- Handler tests for:
  - request validation
  - status code mapping
  - malformed payload handling

### Integration Tests

- Backend integration tests against SQLite database file created per test run.
- End-to-end backend scenarios:
  - create board from template
  - create card
  - move card
  - reload board state
  - delete board and verify cascade
- Frontend integration tests:
  - board detail fetch and render
  - create card flow
  - move card flow with API-backed refresh
  - error state rendering when API fails

### Containerized Test Execution

- Backend tests must run inside the backend container target.
- Frontend tests must run inside the frontend container target.
- Compose must support a reproducible test command for local use and ACP runtimes.
- Containerized execution is the default verification path when running under Claude, Codex, or similar automated environments.

## Development Sequencing

### Build Order

1. Initialize backend project structure, HTTP router, migrations, SQLite connection, and backend Dockerfile - no dependencies
2. Implement board, column, and card schema plus migration runner - depends on step 1
3. Implement repository layer for boards, columns, and cards - depends on step 2
4. Implement service layer with ordering and movement rules - depends on step 3
5. Implement REST handlers and request/response DTOs - depends on step 4
6. Implement board templates and board creation flow - depends on steps 4 and 5
7. Initialize Next.js frontend, API client layer, and frontend Dockerfile - depends on step 5
8. Add Docker Compose orchestration for frontend and backend local startup - depends on steps 1 and 7
9. Build board list and board detail screens - depends on steps 7 and 8
10. Build card and column CRUD UI - depends on steps 7, 8, and 9
11. Build drag-and-drop flow with server-confirmed updates - depends on steps 5, 9, and 10
12. Add automated tests across backend and frontend critical flows - depends on steps 3 through 11
13. Add containerized test commands and compose test profile - depends on steps 8 and 12

### Technical Dependencies

- Go HTTP stack selection and project bootstrap
- Next.js app bootstrap
- SQLite driver and migration tool selection
- Docker and Docker Compose availability in local and agent environments
- Shared local development conventions for frontend-backend base URL and CORS

## Monitoring and Observability

- Structured backend logs for:
  - request method
  - request path
  - status code
  - latency
  - entity identifiers when available
- Error logs for repository failures and invalid move operations
- Frontend error reporting to browser console in MVP
- Container logs must be readable through Compose for frontend and backend services
- No alerting system required for MVP local usage

## Technical Considerations

### Key Decisions

- Decision: Use separated frontend and backend applications.
  - Rationale: Keeps Go domain logic independent and leaves the frontend free to focus on UI concerns.
  - Trade-offs: Requires explicit API contracts and local cross-origin setup.
  - Alternatives rejected: Next.js-only BFF and fully monolithic server-rendered app.

- Decision: Use SQLite for MVP persistence.
  - Rationale: Minimal setup and good fit for a single-user application.
  - Trade-offs: Concurrency and future scaling are limited compared with PostgreSQL.
  - Alternatives rejected: PostgreSQL from day one and in-memory persistence.

- Decision: Use handlers + services + repositories in the Go backend.
  - Rationale: Enough structure to isolate business rules and prepare database migration later.
  - Trade-offs: Slightly more boilerplate than a handler-only MVP.
  - Alternatives rejected: all-in-handlers design and early hexagonal architecture.

- Decision: Avoid optimistic UI updates in MVP.
  - Rationale: Reduces reconciliation complexity and keeps data consistency simple.
  - Trade-offs: Drag-and-drop may feel less immediate than a richer client-side state model.
  - Alternatives rejected: fully optimistic DnD and hybrid optimistic-only move operations.

- Decision: Standardize development and test execution with Docker and Docker Compose.
  - Rationale: Reduces environment drift and makes execution reproducible across local machines and agent runtimes.
  - Trade-offs: Adds container maintenance overhead and a small amount of startup complexity.
  - Alternatives rejected: host-only setup instructions and ad hoc per-agent scripts.

### Known Risks

- Card and column ordering logic can become error-prone during move and delete operations.
  - Mitigation: keep ordering logic centralized in services and test transactional scenarios thoroughly.

- SQLite may show locking issues if multiple tabs write at the same time.
  - Mitigation: enable WAL mode, keep write transactions short, and accept limited MVP concurrency.

- Waiting for server confirmation on drag-and-drop may make interactions feel slower.
  - Mitigation: keep move endpoint fast and provide clear loading feedback during move operations.

- Frontend-backend separation adds local setup overhead.
  - Mitigation: standardize local ports, base URLs, and developer scripts early.

- Containerized development can become slower on some hosts with bind mounts.
  - Mitigation: keep images minimal, scope mounted paths carefully, and allow direct local execution as a fallback for developers when needed.

## Architecture Decision Records

- [ADR-001: MVP Lean for Go Kanban](adrs/adr-001.md) — Selected a lean MVP focused on personal boards, templates, basic cards, drag-and-drop, and durable persistence.
- [ADR-002: Separate Next.js Frontend and Go REST Backend](adrs/adr-002.md) — Chosen to keep UI delivery separate from domain and persistence logic.
- [ADR-003: SQLite for MVP Persistence with Future PostgreSQL Migration Path](adrs/adr-003.md) — Chosen to reduce operational complexity while preserving a clean migration path.
- [ADR-004: Layered Go Backend with Server-Confirmed UI Updates](adrs/adr-004.md) — Chosen to centralize business rules in services and avoid optimistic state complexity in the MVP.
- [ADR-005: Docker and Docker Compose for Reproducible Local and Agent Execution](adrs/adr-005.md) — Chosen to standardize development and test workflows across local machines and ACP runtimes.
