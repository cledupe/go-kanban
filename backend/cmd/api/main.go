package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cledupe/go-kanban/backend/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := realMain(ctx, os.Getenv, app.Run); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func realMain(ctx context.Context, getenv func(string) string, run func(context.Context, app.Config) error) error {
	cfg, err := app.LoadConfig(getenv)
	if err != nil {
		return err
	}

	return run(ctx, cfg)
}
