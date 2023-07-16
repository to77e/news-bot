// Package repository provides work with database.
package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/to77e/news-bot/internal/models"
	"time"
)

// ErrorSourceNotFound - source not found error.
var (
	ErrorSourceNotFound = errors.New("source not found")
)

type dbSource struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	FeedURL     string    `db:"feed_url"`
	CreatedDate time.Time `db:"created_at"`
}

// SourceRepository - source repository.
type SourceRepository struct {
	db *pgx.Conn
}

// Sources - returns all sources.
func (s *SourceRepository) Sources(ctx context.Context) ([]*models.Source, error) {
	const (
		query = `SELECT id, name, feed_url, created_at FROM sources;`
	)

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select sources: %w", err)
	}
	defer rows.Close()

	var sources []*models.Source
	for rows.Next() {
		var source dbSource
		if err := rows.Scan(&source.ID, &source.Name, &source.FeedURL, &source.CreatedDate); err != nil {
			return nil, err
		}

		sources = append(sources, &models.Source{
			ID:          source.ID,
			Name:        source.Name,
			FeedURL:     source.FeedURL,
			CreatedDate: source.CreatedDate,
		})
	}

	return sources, nil
}

// SourceByID - returns source by id.
func (s *SourceRepository) SourceByID(ctx context.Context, id int64) (*models.Source, error) {
	const (
		query = `SELECT id, name, feed_url, created_at FROM sources WHERE id = $1;`
	)

	var source dbSource
	err := s.db.QueryRow(ctx, query, id).Scan(&source.ID, &source.Name, &source.FeedURL, &source.CreatedDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrorSourceNotFound
		}
		return nil, fmt.Errorf("select source by id %d: %w", id, err)
	}

	return (*models.Source)(&source), nil
}

// Add - adds new source.
func (s *SourceRepository) Add(ctx context.Context, source models.Source) (int64, error) {
	const (
		query = `INSERT INTO sources (name, feed_url) VALUES ($1, $2) RETURNING id;`
	)

	var id int64
	err := s.db.QueryRow(ctx, query, source.Name, source.FeedURL).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert source: %w", err)
	}

	return id, nil
}

// Delete - deletes source by id.
func (s *SourceRepository) Delete(ctx context.Context, id int64) error {
	const (
		query = `DELETE FROM sources WHERE id = $1;`
	)

	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete source: %w", err)
	}

	return nil
}
