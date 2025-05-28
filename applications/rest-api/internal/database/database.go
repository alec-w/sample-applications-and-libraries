package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alec-w/sample-applications-and-libraries/applications/rest-api/internal/models"
	"github.com/alec-w/sample-applications-and-libraries/libraries/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	db     *sql.DB
	logger logging.Logger
}

func NewDatabase(ctx context.Context, host, user, dbname, password string, port int, logger logging.Logger) (*Database, error) {
	cfg, err := pgx.ParseConfig(fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", host, port, user, dbname, password))
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}
	return &Database{db: stdlib.OpenDB(*cfg), logger: logger}, nil
}

func (d *Database) ListPosts() ([]models.Post, error) {
	d.logger.Debug("listing posts in database")
	return []models.Post{
		{
			Id:        "one",
			Title:     "test",
			Content:   "Test Post",
			CreatedAt: time.Now(),
		},
	}, nil
}
