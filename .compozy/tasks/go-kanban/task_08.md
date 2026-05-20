---
status: pending
title: Drag-and-drop com persistência via API e cobertura mínima
type: frontend
complexity: high
dependencies:
  - task_07
---

# Task 08: Drag-and-drop com persistência via API e cobertura mínima

## Overview
Implement the drag-and-drop board interaction for moving cards within and across columns, persisting the final order through the backend API. This task completes the MVP's central interaction loop and must preserve data integrity while staying within the server-confirmed update model.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The frontend MUST support moving cards within the same column and across columns through drag-and-drop.
- Drag-and-drop completion MUST persist through the backend move endpoint defined in TechSpec "API Endpoints".
- Visible board state MUST only finalize after the backend confirms the move result.
- The UI MUST provide clear interaction feedback for active drag, drop targets, and move failures.
- The task MUST include automated tests for reorder, cross-column move, failed persistence, and post-move board refresh behavior.
</requirements>

## Subtasks
- [ ] 8.1 Add drag-and-drop interaction to cards in the board detail UI.
- [ ] 8.2 Support card reordering inside a column and movement across columns.
- [ ] 8.3 Persist completed moves through the backend API and refresh board state after success.
- [ ] 8.4 Add clear interaction and failure feedback during drag-and-drop operations.
- [ ] 8.5 Add automated checks covering move success, reordering, cross-column transfer, and failed persistence.

## Implementation Details
Use PRD "Drag-and-Drop", TechSpec "API Endpoints", "Testing Approach", and "Technical Considerations", plus the earlier frontend mutation patterns from Task 07. Preserve the no-optimistic-update rule while still providing clear interaction feedback.

### Relevant Files
- `.compozy/tasks/go-kanban/_prd.md` — Defines drag-and-drop as a core MVP capability and success metric.
- `.compozy/tasks/go-kanban/_techspec.md` — Defines the move endpoint, server-confirmed behavior, and test expectations.
- `.compozy/tasks/go-kanban/adrs/adr-001.md` — Keeps drag-and-drop in the lean MVP scope.
- `.compozy/tasks/go-kanban/adrs/adr-004.md` — Requires server-confirmed UI updates instead of optimistic final state.
- `.compozy/tasks/go-kanban/adrs/adr-005.md` — Requires reproducible test execution for this interaction-heavy flow.

### Dependent Files
- `frontend/components/board/` — Expected home for card drag-and-drop interactions and visual state.
- `frontend/lib/api/` — Existing client boundary for move requests and board refresh behavior.
- `frontend/app/` or `frontend/src/app/` — Likely route-level coordination point for move persistence and rerendering.
- `frontend/tests/` — Expected location for drag-and-drop interaction and integration tests.
- `docker-compose.yml` — May need test command or profile updates for reproducible frontend interaction verification.

### Related ADRs
- [ADR-001: MVP Lean for Go Kanban](adrs/adr-001.md) — Includes drag-and-drop in the MVP success path.
- [ADR-004: Layered Go Backend with Server-Confirmed UI Updates](adrs/adr-004.md) — Requires move completion to be confirmed by the backend.
- [ADR-005: Docker and Docker Compose for Reproducible Local and Agent Execution](adrs/adr-005.md) — Supports consistent verification in local and automated environments.

## Deliverables
- Drag-and-drop interaction for card reordering and cross-column movement
- Backend-backed move persistence with post-success board refresh
- Visual interaction and failure feedback for move operations
- Container-friendly automated verification path for drag-and-drop behavior
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for drag-and-drop persistence flows **(REQUIRED)**

## Tests
- Unit tests:
  - [ ] Drag completion within the same column triggers the correct move request and only finalizes UI after success.
  - [ ] Cross-column drag completion triggers the correct target column and position payload.
  - [ ] Failed move persistence preserves the last confirmed board state and renders error feedback.
- Integration tests:
  - [ ] A user can drag a card within a column, persist the new order, and see the updated board state after refresh.
  - [ ] A user can drag a card to another column, persist the move, and see the updated board state after refresh.
  - [ ] A backend move failure leaves the board in the previously confirmed state and surfaces an error message.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- A user can move cards within and across columns through drag-and-drop with persisted server state
- The board remains consistent after successful and failed move operations
