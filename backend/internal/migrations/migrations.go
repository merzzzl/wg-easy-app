package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var files embed.FS

func Run(ctx context.Context, db *sql.DB) error {
	goose.SetBaseFS(files)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, "sql"); err != nil {
		return fmt.Errorf("run goose migrations: %w", err)
	}

	return nil
}
