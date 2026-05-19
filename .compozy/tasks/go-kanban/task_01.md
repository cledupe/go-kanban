---
status: completed
title: Bootstrap do backend Go com Docker e Compose
type: backend
complexity: high
dependencies: []
---

# Task 01: Bootstrap do backend Go com Docker e Compose

## Overview
Create the initial Go backend structure for the MVP and make it runnable through Docker and Docker Compose from day one. This task establishes the execution baseline required by the architecture, containerization, and layered backend ADRs.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The backend MUST be initialized as a standalone Go application consistent with the TechSpec "System Architecture" section.
- The backend MUST expose a minimal health or readiness endpoint to support local orchestration checks.
- The repository MUST include a backend Dockerfile and a root Docker Compose configuration aligned with TechSpec "Containerization".
- The compose setup MUST allow the backend service to start with mounted source code and stable local ports.
- The task MUST include automated tests for startup, routing baseline, and containerized execution smoke checks.
</requirements>

## Subtasks
- [x] 1.1 Create the backend project structure, module definition, and application entrypoint.
- [x] 1.2 Add the initial HTTP router and a minimal readiness endpoint for orchestration.
- [x] 1.3 Add backend container definitions for development and test execution targets.
- [x] 1.4 Add root Compose orchestration that can start the backend service consistently in local and agent environments.
- [x] 1.5 Add baseline automated checks covering process startup and containerized smoke execution.

## Implementation Details
Create the initial backend folder, Go module, app entrypoint, router boundary, and containerization assets. See TechSpec "System Architecture", "Containerization", and "Development Sequencing" sections for the expected boundaries and startup order.

### Relevant Files
- `.compozy/tasks/go-kanban/_techspec.md` — Defines backend boundaries, containerization requirements, and build order.
- `.compozy/tasks/go-kanban/adrs/adr-002.md` — Confirms the separated frontend and backend architecture.
- `.compozy/tasks/go-kanban/adrs/adr-004.md` — Confirms the layered Go backend approach.
- `.compozy/tasks/go-kanban/adrs/adr-005.md` — Requires Docker and Compose for reproducible execution.

### Dependent Files
- `backend/go.mod` — Required to declare the Go module and dependencies.
- `backend/cmd/api/main.go` — Expected backend application entrypoint.
- `backend/internal/http/` — Expected HTTP boundary for routing and handlers.
- `backend/Dockerfile` — Required for containerized development and test execution.
- `docker-compose.yml` — Required to orchestrate backend startup from the repository root.

### Related ADRs
- [ADR-002: Separate Next.js Frontend and Go REST Backend](adrs/adr-002.md) — Defines the two-application topology.
- [ADR-004: Layered Go Backend with Server-Confirmed UI Updates](adrs/adr-004.md) — Defines backend layering expectations.
- [ADR-005: Docker and Docker Compose for Reproducible Local and Agent Execution](adrs/adr-005.md) — Requires standardized containerized execution.

## Deliverables
- Standalone Go backend project scaffold with runnable entrypoint
- Initial router and readiness endpoint
- Backend Dockerfile with development and test targets
- Root Compose file with backend service wiring
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for backend startup and container smoke execution **(REQUIRED)**

## Tests
- Unit tests:
  - [ ] Router initialization returns a handler that serves the readiness endpoint.
  - [ ] Readiness endpoint returns the expected success status and payload.
  - [ ] Backend startup configuration rejects invalid or missing required settings.
- Integration tests:
  - [ ] Backend process starts and serves the readiness endpoint in a local test environment.
  - [ ] Backend container target builds successfully and starts through Compose without manual host setup.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- The repository can start the backend service through Docker Compose with a stable local port
- The backend has a clear entrypoint and HTTP boundary ready for subsequent persistence and API tasks
