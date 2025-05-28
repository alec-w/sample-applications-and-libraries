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

func (d *Database) ListPosts(ctx context.Context) ([]models.Post, error) {
	d.logger.Debug("listing posts in database")
	rows, err := d.db.QueryContext(ctx, "SELECT id, title, content, created_at FROM posts")
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over posts: %w", err)
	}
	return posts, nil
}

func (d *Database) CreatePost(ctx context.Context, title, content string, createdAt time.Time) (models.Post, error) {
	d.logger.Debug("creating post in database")
	post := models.Post{Title: title, Content: content, CreatedAt: createdAt}
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Post{}, fmt.Errorf("failed to start transaction to insert post into database: %w", err)
	}
	defer tx.Rollback()
	// Insert post
	_, err = tx.ExecContext(ctx, "INSERT INTO posts(title, content, created_at) VALUES($1, $2, $3)", post.Title, post.Content, post.CreatedAt)
	if err != nil {
		return models.Post{}, fmt.Errorf("failed to insert post into database: %w", err)
	}
	// Fetch ID
	row := tx.QueryRowContext(ctx, "SELECT MAX(id) FROM posts")
	if err := row.Scan(&post.Id); err != nil {
		return models.Post{}, fmt.Errorf("failed to fetch newly inserted post id from database: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return models.Post{}, fmt.Errorf("faile to commit insertion of post into database: %w", err)
	}
	return post, nil
}
