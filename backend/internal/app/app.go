package app

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"

	httpapi "github.com/cledupe/go-kanban/backend/internal/http"
	"github.com/cledupe/go-kanban/backend/internal/http/handlers"
	"github.com/cledupe/go-kanban/backend/internal/service"
	"github.com/cledupe/go-kanban/backend/internal/storage/sqlite"
)

type App struct {
	server         *http.Server
	db             *sqlite.DB
	listenAndServe func() error
	serve          func(net.Listener) error
	shutdown       func(context.Context) error
}

var newApp = New

func New(cfg Config) (*App, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	db, err := sqlite.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	if err := sqlite.RunMigrations(db); err != nil {
		return nil, err
	}

	boardRepo := sqlite.NewBoardRepository(db)
	columnRepo := sqlite.NewColumnRepository(db)
	cardRepo := sqlite.NewCardRepository(db)

	boardService := service.NewBoardService(boardRepo, columnRepo, cardRepo)
	columnService := service.NewColumnService(boardRepo, columnRepo)
	cardService := service.NewCardService(columnRepo, cardRepo)

	boardHandler := handlers.NewBoardHandler(boardService)
	columnHandler := handlers.NewColumnHandler(columnService)
	cardHandler := handlers.NewCardHandler(cardService)

	handler := httpapi.NewRouter(httpapi.RouterDependencies{
		BoardHandler:  boardHandler,
		ColumnHandler: columnHandler,
		CardHandler:   cardHandler,
	})

	server := &http.Server{
		Addr:    cfg.Address(),
		Handler: handler,
	}

	return &App{
		server:         server,
		db:             db,
		listenAndServe: server.ListenAndServe,
		serve:          server.Serve,
		shutdown:       server.Shutdown,
	}, nil
}

func Run(ctx context.Context, cfg Config) error {
	application, err := newApp(cfg)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- application.listenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("shutting down...")
		if application.db != nil {
			if err := application.db.Close(); err != nil {
				log.Printf("close db: %v", err)
			}
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return application.shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func (a *App) Serve(listener net.Listener) error {
	return a.serve(listener)
}

func (a *App) Handler() http.Handler {
	return a.server.Handler
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.shutdown(ctx)
}

func (a *App) DB() *sqlite.DB {
	return a.db
}