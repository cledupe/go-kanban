---
status: pending
title: CRUD de columns e cards com atualização confirmada pelo servidor
type: frontend
complexity: high
dependencies:
  - task_06
---

# Task 07: CRUD de columns e cards com atualização confirmada pelo servidor

## Overview
Add the core interactive UI flows for creating, editing, and deleting columns and cards, using server-confirmed updates rather than optimistic local state. This task delivers the main board management capabilities required by the MVP while staying aligned with the approved UI consistency trade-off.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The frontend MUST support create, edit, and delete flows for cards and columns covered by the MVP scope.
- UI state MUST update only after successful backend responses, consistent with ADR-004.
- Validation and error feedback MUST be visible for failed create, update, or delete actions.
- The board detail UI MUST refresh or reconcile server state after each successful mutation.
- The task MUST include automated tests for mutation success paths, validation failures, and server-error handling.
</requirements>

## Subtasks
- [ ] 7.1 Add column create, rename, and delete user flows in the board detail view.
- [ ] 7.2 Add card create, edit, and delete user flows in the board detail view.
- [ ] 7.3 Ensure the UI refreshes server state only after successful mutation responses.
- [ ] 7.4 Add visible validation and failure handling for mutation errors.
- [ ] 7.5 Add automated checks covering successful and failed mutation flows.

## Implementation Details
Follow PRD "Column Management" and "Card Management" plus TechSpec "API Endpoints" and "Technical Considerations". Keep mutation flows aligned with the server-confirmed state model and avoid optimistic local writes.

### Relevant Files
- `.compozy/tasks/go-kanban/_prd.md` — Defines the CRUD scope for columns and cards.
- `.compozy/tasks/go-kanban/_techspec.md` — Defines API endpoints and server-confirmed update behavior.
- `.compozy/tasks/go-kanban/adrs/adr-004.md` — Requires updates to be confirmed by the backend before becoming visible.

### Dependent Files
- `frontend/components/board/` — Expected location for column and card interaction components.
- `frontend/lib/api/` — Existing API client layer used for mutations.
- `frontend/app/` or `frontend/src/app/` — Likely route-level integration point for state refresh behavior.
- `frontend/tests/` — Expected location for mutation and UI-state tests.

### Related ADRs
- [ADR-004: Layered Go Backend with Server-Confirmed UI Updates](adrs/adr-004.md) — Defines the no-optimistic-update interaction model.

## Deliverables
- Column CRUD flows in the board detail UI
- Card CRUD flows in the board detail UI
- Server-confirmed refresh or reconciliation after successful mutations
- Error and validation feedback for failed mutations
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for card and column mutation flows **(REQUIRED)**

## Tests
- Unit tests:
  - [ ] Creating a card or column only updates visible UI after the mutation request succeeds.
  - [ ] Editing a card title or description re-renders the board detail with the server-confirmed values.
  - [ ] Failed create, update, or delete requests render validation or error feedback without corrupting board state.
- Integration tests:
  - [ ] A user can create, edit, and delete cards through the board detail screen with successful server responses.
  - [ ] A user can create, rename, and delete columns through the board detail screen with expected post-mutation refresh.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- A user can manage cards and columns through the frontend without relying on optimistic state
- Failed mutations do not leave the board UI in an inconsistent state
