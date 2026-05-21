---
status: completed
title: Serviços e API REST para boards, columns e cards
type: backend
complexity: high
dependencies:
  - task_02
---

# Task 03: Serviços e API REST para boards, columns e cards

## Overview
Implement the core backend service layer and REST API for the board, column, and card resources. This task exposes the MVP contract consumed by the frontend and centralizes domain rules outside the HTTP and storage layers.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The backend MUST expose the REST endpoints defined in TechSpec "API Endpoints" for boards, columns, and cards.
- Business rules MUST be implemented in services, not in handlers or repositories, consistent with ADR-004.
- HTTP responses MUST map validation, not-found, and conflict conditions to stable status codes.
- Board detail responses MUST return ordered columns and cards suitable for direct frontend rendering.
- The task MUST include automated tests for service behavior, request validation, and API status code mapping.
</requirements>

## Subtasks
- [x] 3.1 Define service interfaces and domain errors for board, column, and card operations.
- [x] 3.2 Implement service behavior for CRUD flows and board detail aggregation.
- [x] 3.3 Add HTTP handlers and request-response DTOs for all MVP resource endpoints.
- [x] 3.4 Wire services and handlers into the backend application startup and routing tree.
- [x] 3.5 Add automated checks covering service rules, validation failures, and endpoint contracts.

## Implementation Details
Use TechSpec "Core Interfaces", "API Endpoints", and "Testing Approach" as the source of truth for this task. Keep transport mapping in handlers, business rules in services, and persistence access in repositories.

### Relevant Files
- `.compozy/tasks/go-kanban/_techspec.md` — Defines service contracts, domain entities, and REST surface.
- `.compozy/tasks/go-kanban/adrs/adr-002.md` — Requires a clean API boundary between frontend and backend.
- `.compozy/tasks/go-kanban/adrs/adr-004.md` — Requires the layered backend structure and server-confirmed state flow.

### Dependent Files
- `backend/internal/service/` — Expected location for business services and domain orchestration.
- `backend/internal/http/handlers/` — Expected location for endpoint handlers and DTO mapping.
- `backend/internal/http/router.go` — Expected route registration point.
- `backend/internal/domain/` — Likely location for domain errors and interfaces reused by handlers and services.
- `backend/internal/storage/sqlite/` — Repository implementations consumed by services.

### Related ADRs
- [ADR-002: Separate Next.js Frontend and Go REST Backend](adrs/adr-002.md) — Defines the explicit HTTP contract boundary.
- [ADR-004: Layered Go Backend with Server-Confirmed UI Updates](adrs/adr-004.md) — Requires handlers, services, and repositories with clear responsibilities.

## Deliverables
- Service layer for boards, columns, and cards
- REST handlers and request-response DTOs for MVP endpoints
- Stable error-to-status-code mapping for API consumers
- Board detail endpoint suitable for frontend rendering
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for REST API resource flows **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Service create and update operations reject invalid input with domain errors.
  - [x] Board detail service returns columns and cards in the expected order.
  - [x] Handler validation maps malformed payloads and unknown resource identifiers to correct HTTP responses.
- Integration tests:
  - [x] `POST /api/boards` creates a board and returns the expected response contract.
  - [x] `GET /api/boards/:id` returns the full board shape with ordered columns and cards.
  - [x] `PATCH` and `DELETE` resource endpoints return stable status codes for success and missing records.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- The frontend has a complete REST API contract for board, column, and card management
- Domain rules are centralized in services instead of being split across handlers and repositories
