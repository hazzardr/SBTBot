package sbtb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/hazzardr/sbtbot/internal/generated"
)

type DB struct {
	path    string
	queries *generated.Queries
}

func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	queries := generated.New(db)
	return &DB{path: path, queries: queries}, nil
}

func (db *DB) ListGenres(ctx context.Context) ([]generated.Genre, error) {
	genres, err := db.queries.ListGenres(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve genres: %w", err)
	}
	return genres, nil
}

func (db *DB) AddGenre(ctx context.Context, name string) error {
	g, err := db.queries.AddGenre(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to add genre: %w", err)
	}
	slog.Info("added genre:", g.Name)
	return nil
}
