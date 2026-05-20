---
status: pending
title: UI de lista de boards e detalhe do board
type: frontend
complexity: high
dependencies:
  - task_05
---

# Task 06: UI de lista de boards e detalhe do board

## Overview
Build the first usable frontend screens: the board list and the board detail view that renders ordered columns and cards from the backend API. This task delivers the baseline navigation and visibility required for all later board interactions.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The frontend MUST render a board list view and a board detail view aligned with the PRD user flow.
- The board detail view MUST render ordered columns and cards using the backend response contract from Task 03.
- The UI MUST handle loading, empty, and API error states explicitly.
- The implementation MUST support responsive rendering consistent with the TechSpec and PRD UX guidance.
- The task MUST include automated tests for board list rendering, board detail rendering, and API error handling states.
</requirements>

## Subtasks
- [ ] 6.1 Add board list data loading and rendering for the single-user workspace.
- [ ] 6.2 Add board detail routing and data loading based on the selected board identifier.
- [ ] 6.3 Render ordered columns and cards in the board detail view using backend data.
- [ ] 6.4 Add explicit loading, empty, and error states for list and detail screens.
- [ ] 6.5 Add automated checks covering main render paths and state transitions.

## Implementation Details
Use PRD "First-time user flow" and "Regular use flow" together with TechSpec "System Architecture" and "Testing Approach". Keep state management simple and centered on server-fetched data.

### Relevant Files
- `.compozy/tasks/go-kanban/_prd.md` — Defines the board list and board detail user journeys.
- `.compozy/tasks/go-kanban/_techspec.md` — Defines the board detail payload shape and frontend responsibilities.
- `.compozy/tasks/go-kanban/adrs/adr-001.md` — Keeps the MVP focused on functional usability.
- `.compozy/tasks/go-kanban/adrs/adr-002.md` — Reinforces frontend consumption of the backend API contract.

### Dependent Files
- `frontend/app/` or `frontend/src/app/` — Expected route and screen implementation area.
- `frontend/lib/api/` — Existing API client layer consumed by the screens.
- `frontend/components/board/` — Likely home for board list, board detail, column, and card rendering primitives.
- `frontend/tests/` — Expected location for UI and integration-oriented frontend tests.

### Related ADRs
- [ADR-001: MVP Lean for Go Kanban](adrs/adr-001.md) — Requires the first usable board experience without extra scope.
- [ADR-002: Separate Next.js Frontend and Go REST Backend](adrs/adr-002.md) — Requires board views to consume the backend API contract.

## Deliverables
- Board list screen for existing boards
- Board detail screen with ordered columns and cards
- Explicit loading, empty, and error states for the main board flows
- Responsive baseline layout for board navigation and rendering
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for board list and board detail rendering **(REQUIRED)**

## Tests
- Unit tests:
  - [ ] Board list view renders board names returned by the API client.
  - [ ] Board detail view renders columns and cards in the order returned by the backend.
  - [ ] Loading, empty, and error states render the expected UI for list and detail flows.
- Integration tests:
  - [ ] Navigating from the board list to a board detail view fetches and renders the selected board.
  - [ ] API failure on board list or board detail load renders the expected error state without crashing the app.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- A user can see available boards and open a board detail screen backed by live API data
- Ordered columns and cards render correctly with explicit loading and error handling
