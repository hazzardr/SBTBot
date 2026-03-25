package sbtb

import (
	"context"
	"database/sql"
	"fmt"

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
