---
status: completed
title: Persistência SQLite, migrações e repositórios base
type: backend
complexity: high
dependencies:
  - task_01
---

# Task 02: Persistência SQLite, migrações e repositórios base

## Overview
Add the durable persistence foundation for boards, columns, and cards using SQLite, explicit migrations, and repository contracts. This task turns the backend scaffold into a stateful service while preserving a clean path for future PostgreSQL migration.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The backend MUST persist boards, columns, and cards in SQLite according to TechSpec "Data Models".
- Schema changes MUST be represented through explicit migrations rather than implicit startup creation.
- Repository interfaces and implementations MUST support future database replacement without moving business rules out of services.
- Persistence setup MUST configure SQLite for durable local use and short write transactions.
- The task MUST include automated tests for migrations, CRUD baseline behavior, ordering fields, and cascade deletion rules.
</requirements>

## Subtasks
- [x] 2.1 Add SQLite connection management and migration execution to the backend startup flow.
- [x] 2.2 Create the initial schema for boards, columns, and cards with required foreign keys and ordering fields.
- [x] 2.3 Define repository contracts for boards, columns, and cards in the backend domain boundary.
- [x] 2.4 Implement SQLite-backed repositories for core create, read, update, and delete behavior.
- [x] 2.5 Add automated checks for schema creation, repository behavior, and cascade deletion semantics.

## Implementation Details
Implement the persistence layer described in TechSpec "Data Models", "Testing Approach", and "Technical Considerations". Repository code should remain storage-focused, leaving ordering and workflow rules to later service-layer tasks.

### Relevant Files
- `.compozy/tasks/go-kanban/_techspec.md` — Defines SQLite schema, repository boundaries, and migration expectations.
- `.compozy/tasks/go-kanban/adrs/adr-003.md` — Defines SQLite as the MVP persistence engine with a migration path.
- `.compozy/tasks/go-kanban/adrs/adr-004.md` — Defines repository separation from services.

### Dependent Files
- `backend/internal/storage/sqlite/` — Expected SQLite connection and repository implementation area.
- `backend/internal/domain/` — Expected home for entity and repository contract definitions.
- `backend/migrations/` — Expected location for SQL migration files.
- `backend/internal/app/` — Likely startup integration point for database initialization.
- `backend/Dockerfile` — May require updates to support migration and persistence execution in containers.

### Related ADRs
- [ADR-003: SQLite for MVP Persistence with Future PostgreSQL Migration Path](adrs/adr-003.md) — Defines the persistence choice and migration boundary.
- [ADR-004: Layered Go Backend with Server-Confirmed UI Updates](adrs/adr-004.md) — Requires business logic to stay out of repositories.
- [ADR-005: Docker and Docker Compose for Reproducible Local and Agent Execution](adrs/adr-005.md) — Ensures persistence works consistently in containers.

## Deliverables
- SQLite connection and migration execution integrated into backend startup
- Initial migration files for boards, columns, and cards
- Repository contracts and SQLite implementations for base persistence operations
- Container-compatible persistence configuration
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for migrations and repository-backed persistence flows **(REQUIRED)**

## Tests
- Unit tests:
  - [ ] Migration runner applies the initial schema to an empty SQLite database.
  - [ ] Board, column, and card repositories persist and load records with required fields populated.
  - [ ] Repository operations preserve position fields and required foreign key relationships.
- Integration tests:
  - [ ] Backend startup applies migrations and opens a usable SQLite database in a clean environment.
  - [ ] Deleting a board removes dependent columns and cards through cascade behavior.
  - [ ] Containerized backend execution can persist and reload data across process restarts using the configured volume or data path.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- The backend can create and read persistent board data through SQLite-backed repositories
- Schema creation is repeatable through explicit migrations and works in local and containerized environments
