package httpapi

import (
	"net/http"

	"github.com/cledupe/go-kanban/backend/internal/http/handlers"
)

type RouterDependencies struct {
	BoardHandler   *handlers.BoardHandler
	ColumnHandler  *handlers.ColumnHandler
	CardHandler    *handlers.CardHandler
}

func NewRouter(deps RouterDependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /readyz", readinessHandler)

	mux.HandleFunc("GET /api/boards", deps.BoardHandler.List)
	mux.HandleFunc("POST /api/boards", deps.BoardHandler.Create)
	mux.HandleFunc("GET /api/boards/{id}", deps.BoardHandler.Get)
	mux.HandleFunc("PATCH /api/boards/{id}", deps.BoardHandler.Update)
	mux.HandleFunc("DELETE /api/boards/{id}", deps.BoardHandler.Delete)

	mux.HandleFunc("POST /api/boards/{boardId}/columns", deps.ColumnHandler.Create)
	mux.HandleFunc("PATCH /api/columns/{id}", deps.ColumnHandler.Update)
	mux.HandleFunc("DELETE /api/columns/{id}", deps.ColumnHandler.Delete)

	mux.HandleFunc("POST /api/columns/{columnId}/cards", deps.CardHandler.Create)
	mux.HandleFunc("PATCH /api/cards/{id}", deps.CardHandler.Update)
	mux.HandleFunc("POST /api/cards/{id}/move", deps.CardHandler.Move)
	mux.HandleFunc("DELETE /api/cards/{id}", deps.CardHandler.Delete)

	return mux
}