---
status: pending
title: Bootstrap do frontend Next.js com cliente HTTP e Compose
type: frontend
complexity: high
dependencies:
  - task_03
---

# Task 05: Bootstrap do frontend Next.js com cliente HTTP e Compose

## Overview
Create the standalone Next.js frontend foundation and wire it to the backend REST API through a stable HTTP client boundary. This task establishes the UI runtime and containerized workflow required before feature screens can be built.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The frontend MUST be initialized as a standalone Next.js application consistent with ADR-002.
- The frontend MUST consume the Go backend through a dedicated HTTP client boundary rather than calling ad hoc endpoints throughout the UI.
- The repository MUST include frontend container definitions and Compose integration aligned with TechSpec "Containerization".
- The frontend runtime MUST receive backend base URL configuration through environment variables suitable for local and agent execution.
- The task MUST include automated tests for frontend startup, API client configuration, and containerized smoke execution.
</requirements>

## Subtasks
- [ ] 5.1 Create the frontend application structure and initial layout shell.
- [ ] 5.2 Add a dedicated API client layer for backend communication and environment-based configuration.
- [ ] 5.3 Add frontend Dockerfile targets for development and test execution.
- [ ] 5.4 Extend Compose orchestration to run frontend and backend together with stable networking.
- [ ] 5.5 Add automated checks covering frontend startup, configuration loading, and API client baseline behavior.

## Implementation Details
Follow TechSpec "System Architecture", "Containerization", and "Development Sequencing" for the frontend boundary and startup responsibilities. Keep API communication centralized to support later UI tasks and testability.

### Relevant Files
- `.compozy/tasks/go-kanban/_techspec.md` — Defines the separated frontend architecture and container expectations.
- `.compozy/tasks/go-kanban/adrs/adr-002.md` — Requires a standalone Next.js frontend.
- `.compozy/tasks/go-kanban/adrs/adr-005.md` — Requires Docker and Compose support for frontend execution.

### Dependent Files
- `frontend/package.json` — Required to define the frontend runtime and scripts.
- `frontend/app/` or `frontend/src/app/` — Expected application shell area.
- `frontend/lib/api/` — Expected home for backend API client logic.
- `frontend/Dockerfile` — Required for frontend containerized development and test execution.
- `docker-compose.yml` — Must be updated to include frontend orchestration.

### Related ADRs
- [ADR-002: Separate Next.js Frontend and Go REST Backend](adrs/adr-002.md) — Defines the standalone frontend boundary.
- [ADR-005: Docker and Docker Compose for Reproducible Local and Agent Execution](adrs/adr-005.md) — Requires frontend containerization and Compose support.

## Deliverables
- Standalone Next.js frontend scaffold
- Centralized API client configured through environment variables
- Frontend Dockerfile with development and test targets
- Compose orchestration updated to start frontend with backend
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for frontend startup and API client smoke behavior **(REQUIRED)**

## Tests
- Unit tests:
  - [ ] Frontend configuration resolves the backend base URL from the expected environment variables.
  - [ ] API client constructs requests using the configured backend base URL and expected resource prefixes.
  - [ ] Frontend shell renders the baseline application frame without runtime errors.
- Integration tests:
  - [ ] Frontend container target builds successfully and starts in Compose alongside the backend.
  - [ ] Frontend can reach the backend readiness or baseline API route through the configured service topology.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- The repository can start frontend and backend together through Compose
- The frontend has a stable HTTP client boundary ready for board and card UI implementation
