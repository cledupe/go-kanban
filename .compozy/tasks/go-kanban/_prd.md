# PRD: Go Kanban

## Overview

Go Kanban is a personal Kanban board application that helps individuals visualize and manage their workflow. It solves the problem of task overload and lack of workflow visibility for solo workers, freelancers, and developers who need a lightweight alternative to heavy project management tools.

The product is for individuals who want a simple, fast, and self-hosted Kanban board without the complexity of team collaboration features. It provides immediate value through pre-built templates while allowing full customization as the user's workflow evolves.

## Goals

- Deliver a working Kanban board with board creation, card management, and drag-and-drop within a single development cycle
- Support at least 3 pre-built templates out of the box (Basic Kanban, Bug Tracker, Content Pipeline)
- Enable users to create, customize, and delete boards and columns freely
- Persist all data reliably so no work is lost between sessions
- Success criteria: A user can create a board from a template, add cards, move them between columns via drag-and-drop, and see their changes persisted on reload

## User Stories

**Primary Persona: Solo Worker (developer, freelancer, creator)**

- As a solo worker, I want to create a Kanban board from a template so that I can start managing my tasks immediately without setup overhead
- As a user, I want to add cards with a title and description so that I can capture my work items
- As a user, I want to drag cards between columns so that I can update task status visually
- As a user, I want to create custom boards with my own columns so that I can model my unique workflow
- As a user, I want to edit and delete cards so that I can keep my board accurate
- As a user, I want to create, rename, and delete columns so that I can adapt the board as my process changes
- As a user, I want to delete boards I no longer need so that my workspace stays organized

## Core Features

### Board Management
- Create boards from pre-defined templates (Basic Kanban: To Do / In Progress / Done; Bug Tracker: Backlog / Investigating / Fixing / Verified; Content Pipeline: Ideas / Drafting / Review / Published)
- Create blank boards with default columns
- Rename and delete boards
- Board list view showing all user boards

### Column Management
- Add, rename, reorder, and delete columns within a board
- Columns define the workflow stages
- Visual column headers with card count

### Card Management
- Create cards with title (required) and description (optional)
- Edit card title and description inline or in a detail view
- Delete cards
- Cards are ordered within columns and can be reordered

### Drag-and-Drop
- Drag cards between columns to change status
- Drag cards within a column to reorder
- Visual feedback during drag operation (drop zones, card ghost)
- Changes persist immediately on drop

### Template System
- Pre-built templates ship with the application
- Templates define initial board name and column structure
- Users can customize templates' boards after creation like any other board

## User Experience

**First-time user flow:**
1. User lands on the app and sees a welcome screen with option to create a board
2. User chooses "Create from template" or "Create blank board"
3. If template: user selects from available templates and confirms
4. Board appears with pre-configured columns
5. User clicks "Add card" in any column to create their first card

**Regular use flow:**
1. User opens the app and sees their board list
2. User selects a board
3. User manages cards: creates, edits, moves via drag-and-drop, deletes
4. Changes are saved automatically

**UI/UX considerations:**
- Clean, minimal interface focused on the board view
- Board view is the primary screen; board list is secondary
- Drag-and-drop should have clear visual indicators (highlighted drop zones)
- Card detail can be a modal or inline expansion
- Responsive layout: board scrolls horizontally on smaller screens
- No authentication UI in MVP (single user, no login)

## High-Level Technical Constraints

- The application must persist data across server restarts
- Board rendering should feel responsive with up to 100 cards per board
- The system must support concurrent access from a single user (multiple browser tabs)
- Data must be backed by a durable store (not in-memory only)

## Non-Goals (Out of Scope)

- User authentication and multi-user support
- Real-time collaboration or WebSocket sync
- File attachments, checklists, or comments on cards
- Labels, tags, or priority indicators
- Due dates or reminders
- Board sharing or public links
- Analytics, metrics, or reporting
- Mobile native app (responsive web only)
- Search and filter functionality
- Undo/redo actions
- API for third-party integrations

## Phased Rollout Plan

### MVP (Phase 1)
- Board creation from templates and blank
- CRUD for boards, columns, and cards
- Drag-and-drop card movement between columns and reordering
- PostgreSQL persistence
- Single-user, no authentication
- **Success criteria:** User can create a board, manage cards, move them with drag-and-drop, and see changes persist after page reload

### Phase 2
- User authentication (single user with password)
- Card labels/tags and priority indicators
- Due dates on cards
- Board search and filter
- **Success criteria:** Authenticated user can organize cards with tags and due dates, and find cards quickly

### Phase 3
- Card checklists and comments
- Board analytics (cycle time, throughput)
- Card movement history
- Export/import boards
- **Success criteria:** User can track their productivity and export their data

## Success Metrics

- User can create their first board and add a card within 60 seconds of opening the app
- Board loads and renders within 2 seconds for boards with up to 100 cards
- Zero data loss: all card movements and edits persist correctly
- Drag-and-drop interactions complete within 200ms of drop

## Risks and Mitigations

**Adoption risk:** Users may find a single-user Kanban too limiting compared to free tiers of Trello/Notion.
- *Mitigation:* Emphasize simplicity, speed, and self-hosted privacy as differentiators

**Scope creep risk:** Easy to add "just one more feature" (auth, tags, due dates) during MVP development.
- *Mitigation:* Strict adherence to Non-Goals list; defer all additions to Phase 2+

**Template relevance risk:** Pre-built templates may not match user workflows.
- *Mitigation:* Templates are starting points; users can fully customize columns after creation

## Architecture Decision Records

- [ADR-001: MVP Lean for Go Kanban](adrs/adr-001.md) — Selected MVP Lean approach: personal boards with templates, basic cards with drag-and-drop, PostgreSQL persistence, clean functional interface

## Open Questions

- Should the app include a simple landing page explaining the product, or go directly to the board view?
- What happens to cards when a column is deleted — should they move to a default column or be deleted too?
- Should templates be user-editable (save custom templates) in a future phase?
