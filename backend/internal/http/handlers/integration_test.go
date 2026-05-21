package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/cledupe/go-kanban/backend/internal/domain"
	"github.com/cledupe/go-kanban/backend/internal/service"
	"github.com/cledupe/go-kanban/backend/internal/storage/sqlite"
)

func setupIntegration(t *testing.T) http.Handler {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := sqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := sqlite.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	boardRepo := sqlite.NewBoardRepository(db)
	columnRepo := sqlite.NewColumnRepository(db)
	cardRepo := sqlite.NewCardRepository(db)

	boardService := service.NewBoardService(boardRepo, columnRepo, cardRepo)
	columnService := service.NewColumnService(boardRepo, columnRepo)
	cardService := service.NewCardService(columnRepo, cardRepo)

	mux := http.NewServeMux()
	boardHandler := NewBoardHandler(boardService)
	columnHandler := NewColumnHandler(columnService)
	cardHandler := NewCardHandler(cardService)

	mux.HandleFunc("GET /api/boards", boardHandler.List)
	mux.HandleFunc("POST /api/boards", boardHandler.Create)
	mux.HandleFunc("GET /api/boards/{id}", boardHandler.Get)
	mux.HandleFunc("PATCH /api/boards/{id}", boardHandler.Update)
	mux.HandleFunc("DELETE /api/boards/{id}", boardHandler.Delete)

	mux.HandleFunc("POST /api/boards/{boardId}/columns", columnHandler.Create)
	mux.HandleFunc("PATCH /api/columns/{id}", columnHandler.Update)
	mux.HandleFunc("DELETE /api/columns/{id}", columnHandler.Delete)

	mux.HandleFunc("POST /api/columns/{columnId}/cards", cardHandler.Create)
	mux.HandleFunc("PATCH /api/cards/{id}", cardHandler.Update)
	mux.HandleFunc("POST /api/cards/{id}/move", cardHandler.Move)
	mux.HandleFunc("DELETE /api/cards/{id}", cardHandler.Delete)

	return mux
}

func TestIntegrationCreateAndGetBoard(t *testing.T) {
	handler := setupIntegration(t)

	body, _ := json.Marshal(CreateBoardRequest{Name: "My Board"})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create board: expected 201, got %d", rec.Code)
	}

	var board domain.Board
	if err := json.NewDecoder(rec.Body).Decode(&board); err != nil {
		t.Fatalf("decode board: %v", err)
	}
	if board.Name != "My Board" {
		t.Fatalf("expected 'My Board', got %q", board.Name)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/boards/"+board.ID, nil)
	req.SetPathValue("id", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("get board: expected 200, got %d", rec.Code)
	}

	var detail domain.BoardDetail
	if err := json.NewDecoder(rec.Body).Decode(&detail); err != nil {
		t.Fatalf("decode detail: %v", err)
	}
	if detail.Name != "My Board" {
		t.Fatalf("expected 'My Board', got %q", detail.Name)
	}
}

func TestIntegrationCreateBoardFromBasicKanbanTemplate(t *testing.T) {
	handler := setupIntegration(t)

	tmpl := "basic-kanban"
	body, _ := json.Marshal(CreateBoardRequest{Name: "ignored", Template: &tmpl})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create board from template: expected 201, got %d", rec.Code)
	}

	var board domain.Board
	if err := json.NewDecoder(rec.Body).Decode(&board); err != nil {
		t.Fatalf("decode board: %v", err)
	}
	if board.Name != "Basic Kanban" {
		t.Fatalf("expected board name 'Basic Kanban', got %q", board.Name)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/boards/"+board.ID, nil)
	req.SetPathValue("id", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var detail domain.BoardDetail
	json.NewDecoder(rec.Body).Decode(&detail)

	if len(detail.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(detail.Columns))
	}
	expected := []string{"To Do", "In Progress", "Done"}
	for i, col := range detail.Columns {
		if col.Name != expected[i] {
			t.Fatalf("column %d: expected %q, got %q", i, expected[i], col.Name)
		}
		if col.Position != i {
			t.Fatalf("column %d: expected position %d, got %d", i, i, col.Position)
		}
	}
}

func TestIntegrationCreateBoardFromBugTrackerTemplate(t *testing.T) {
	handler := setupIntegration(t)

	tmpl := "bug-tracker"
	body, _ := json.Marshal(CreateBoardRequest{Name: "ignored", Template: &tmpl})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create board from template: expected 201, got %d", rec.Code)
	}

	var board domain.Board
	json.NewDecoder(rec.Body).Decode(&board)
	if board.Name != "Bug Tracker" {
		t.Fatalf("expected board name 'Bug Tracker', got %q", board.Name)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/boards/"+board.ID, nil)
	req.SetPathValue("id", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var detail domain.BoardDetail
	json.NewDecoder(rec.Body).Decode(&detail)

	expected := []string{"Backlog", "Investigating", "Fixing", "Verified"}
	if len(detail.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(detail.Columns))
	}
	for i, col := range detail.Columns {
		if col.Name != expected[i] {
			t.Fatalf("column %d: expected %q, got %q", i, expected[i], col.Name)
		}
		if col.Position != i {
			t.Fatalf("column %d: expected position %d, got %d", i, i, col.Position)
		}
	}
}

func TestIntegrationCreateBoardFromContentPipelineTemplate(t *testing.T) {
	handler := setupIntegration(t)

	tmpl := "content-pipeline"
	body, _ := json.Marshal(CreateBoardRequest{Name: "ignored", Template: &tmpl})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create board from template: expected 201, got %d", rec.Code)
	}

	var board domain.Board
	json.NewDecoder(rec.Body).Decode(&board)
	if board.Name != "Content Pipeline" {
		t.Fatalf("expected board name 'Content Pipeline', got %q", board.Name)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/boards/"+board.ID, nil)
	req.SetPathValue("id", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var detail domain.BoardDetail
	json.NewDecoder(rec.Body).Decode(&detail)

	expected := []string{"Ideas", "Drafting", "Review", "Published"}
	if len(detail.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(detail.Columns))
	}
	for i, col := range detail.Columns {
		if col.Name != expected[i] {
			t.Fatalf("column %d: expected %q, got %q", i, expected[i], col.Name)
		}
		if col.Position != i {
			t.Fatalf("column %d: expected position %d, got %d", i, i, col.Position)
		}
	}
}

func TestIntegrationCreateBoardWithInvalidTemplateReturns400(t *testing.T) {
	handler := setupIntegration(t)

	tmpl := "nonexistent"
	body, _ := json.Marshal(CreateBoardRequest{Name: "Board", Template: &tmpl})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid template, got %d", rec.Code)
	}
}

func TestIntegrationCreateColumnAndCard(t *testing.T) {
	handler := setupIntegration(t)

	body, _ := json.Marshal(CreateBoardRequest{Name: "Board"})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var board domain.Board
	json.NewDecoder(rec.Body).Decode(&board)

	colBody, _ := json.Marshal(CreateColumnRequest{Name: "To Do"})
	req = httptest.NewRequest(http.MethodPost, "/api/boards/"+board.ID+"/columns", bytes.NewReader(colBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("boardId", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create column: expected 201, got %d", rec.Code)
	}

	var col domain.Column
	json.NewDecoder(rec.Body).Decode(&col)
	if col.Name != "To Do" {
		t.Fatalf("expected column name 'To Do', got %q", col.Name)
	}

	cardBody, _ := json.Marshal(CreateCardRequest{Title: "Task 1", Description: "Desc"})
	req = httptest.NewRequest(http.MethodPost, "/api/columns/"+col.ID+"/cards", bytes.NewReader(cardBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("columnId", col.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create card: expected 201, got %d", rec.Code)
	}

	var card domain.Card
	json.NewDecoder(rec.Body).Decode(&card)
	if card.Title != "Task 1" {
		t.Fatalf("expected card title 'Task 1', got %q", card.Title)
	}
}

func TestIntegrationBoardDetailWithOrderedColumnsAndCards(t *testing.T) {
	handler := setupIntegration(t)

	body, _ := json.Marshal(CreateBoardRequest{Name: "Board"})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var board domain.Board
	json.NewDecoder(rec.Body).Decode(&board)

	for _, name := range []string{"To Do", "In Progress", "Done"} {
		colBody, _ := json.Marshal(CreateColumnRequest{Name: name})
		req = httptest.NewRequest(http.MethodPost, "/api/boards/"+board.ID+"/columns", bytes.NewReader(colBody))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("boardId", board.ID)
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/boards/"+board.ID, nil)
	req.SetPathValue("id", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var detail domain.BoardDetail
	json.NewDecoder(rec.Body).Decode(&detail)

	if len(detail.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(detail.Columns))
	}
	if detail.Columns[0].Name != "To Do" {
		t.Fatalf("expected first column 'To Do', got %q", detail.Columns[0].Name)
	}
}

func TestIntegrationPatchAndDeleteReturnStableStatusCodes(t *testing.T) {
	handler := setupIntegration(t)

	body, _ := json.Marshal(CreateBoardRequest{Name: "Board"})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var board domain.Board
	json.NewDecoder(rec.Body).Decode(&board)

	patchBody, _ := json.Marshal(UpdateBoardRequest{Name: "Updated"})
	req = httptest.NewRequest(http.MethodPatch, "/api/boards/"+board.ID, bytes.NewReader(patchBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("patch board: expected 200, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/boards/"+board.ID, nil)
	req.SetPathValue("id", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete board: expected 204, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/boards/"+board.ID, nil)
	req.SetPathValue("id", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("get deleted board: expected 404, got %d", rec.Code)
	}
}

func TestIntegrationMoveCard(t *testing.T) {
	handler := setupIntegration(t)

	body, _ := json.Marshal(CreateBoardRequest{Name: "Board"})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var board domain.Board
	json.NewDecoder(rec.Body).Decode(&board)

	col1Body, _ := json.Marshal(CreateColumnRequest{Name: "To Do"})
	req = httptest.NewRequest(http.MethodPost, "/api/boards/"+board.ID+"/columns", bytes.NewReader(col1Body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("boardId", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	var col1 domain.Column
	json.NewDecoder(rec.Body).Decode(&col1)

	col2Body, _ := json.Marshal(CreateColumnRequest{Name: "Done"})
	req = httptest.NewRequest(http.MethodPost, "/api/boards/"+board.ID+"/columns", bytes.NewReader(col2Body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("boardId", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	var col2 domain.Column
	json.NewDecoder(rec.Body).Decode(&col2)

	cardBody, _ := json.Marshal(CreateCardRequest{Title: "Task"})
	req = httptest.NewRequest(http.MethodPost, "/api/columns/"+col1.ID+"/cards", bytes.NewReader(cardBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("columnId", col1.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	var card domain.Card
	json.NewDecoder(rec.Body).Decode(&card)

	moveBody, _ := json.Marshal(MoveCardRequest{TargetColumnID: col2.ID, Position: 0})
	req = httptest.NewRequest(http.MethodPost, "/api/cards/"+card.ID+"/move", bytes.NewReader(moveBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", card.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("move card: expected 200, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/boards/"+board.ID, nil)
	req.SetPathValue("id", board.ID)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var detail domain.BoardDetail
	json.NewDecoder(rec.Body).Decode(&detail)

	if len(detail.Columns[1].Cards) != 1 {
		t.Fatalf("expected 1 card in 'Done' column after move, got %d", len(detail.Columns[1].Cards))
	}
	if detail.Columns[1].Cards[0].Title != "Task" {
		t.Fatalf("expected card title 'Task', got %q", detail.Columns[1].Cards[0].Title)
	}
}

func TestIntegrationMalformedPayloadReturns400(t *testing.T) {
	handler := setupIntegration(t)

	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for malformed payload, got %d", rec.Code)
	}
}

func TestIntegrationUnknownResourceReturns404(t *testing.T) {
	handler := setupIntegration(t)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown resource, got %d", rec.Code)
	}
}