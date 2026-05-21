---
status: completed
title: Templates de board e regras de criação inicial
type: backend
complexity: medium
dependencies:
  - task_03
---

# Task 04: Templates de board e regras de criação inicial

## Overview
Add the template-backed board creation flow required for the MVP onboarding experience. This task turns generic board creation into a user-facing capability aligned with the PRD success criteria for fast first use.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The backend MUST support creation of blank boards and the three MVP templates defined in the PRD.
- Template creation MUST initialize ordered columns appropriate to the selected template.
- Template resolution MUST be validated and return a stable client error for unknown template identifiers.
- The implementation MUST preserve the same service and API boundaries established by earlier backend tasks.
- The task MUST include automated tests for template selection, resulting column structure, and invalid template handling.
</requirements>

## Subtasks
- [ ] 4.1 Add template definitions for Basic Kanban, Bug Tracker, and Content Pipeline.
- [ ] 4.2 Extend board creation services to support blank and template-backed initialization flows.
- [ ] 4.3 Ensure board creation API contracts accept and validate template identifiers consistently.
- [ ] 4.4 Persist template-created columns with stable ordering.
- [ ] 4.5 Add automated checks covering valid and invalid template creation paths.

## Implementation Details
Use PRD "Template System" and TechSpec "API Endpoints" plus "Technical Considerations" to implement the initial board creation flows. Keep template definitions static and intentionally narrow for the MVP.

### Relevant Files
- `.compozy/tasks/go-kanban/_prd.md` — Defines the three required templates and fast-start user flow.
- `.compozy/tasks/go-kanban/_techspec.md` — Defines template-backed board creation as part of the backend build order.
- `.compozy/tasks/go-kanban/adrs/adr-001.md` — Confirms the lean MVP with immediate value from prebuilt templates.
- `.compozy/tasks/go-kanban/adrs/adr-004.md` — Requires template logic to live in services rather than handlers.

### Dependent Files
- `backend/internal/service/board/` — Likely location for board creation orchestration.
- `backend/internal/http/handlers/boards.go` — Expected board creation request boundary.
- `backend/internal/domain/templates/` — Expected home for static template definitions.
- `backend/internal/storage/sqlite/` — Repositories that persist template-created columns.

### Related ADRs
- [ADR-001: MVP Lean for Go Kanban](adrs/adr-001.md) — Makes templates part of the immediate MVP value.
- [ADR-004: Layered Go Backend with Server-Confirmed UI Updates](adrs/adr-004.md) — Keeps template behavior in service-level orchestration.

## Deliverables
- Static template definitions for the three MVP board presets
- Board creation flow supporting blank and template-backed initialization
- API validation for template identifiers
- Persisted ordered columns created from template selection
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for template-backed board creation **(REQUIRED)**

## Tests
- Unit tests:
  - [ ] Basic Kanban template resolves to exactly three ordered columns.
  - [ ] Bug Tracker and Content Pipeline templates resolve to the expected ordered column sets.
  - [ ] Unknown template identifiers are rejected with the expected domain or validation error.
- Integration tests:
  - [ ] `POST /api/boards` with each supported template creates the board and expected initial columns.
  - [ ] `POST /api/boards` without a template creates a blank board using the defined default behavior.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- A user can create a board from any required MVP template through the backend API
- Template-backed board creation produces the correct initial column structure with stable ordering
