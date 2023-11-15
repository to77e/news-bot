package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/to77e/news-fetching-bot/internal/models"
)

type dbArticle struct {
	ID            int64        `db:"id"`
	SourceID      int64        `db:"source_id"`
	Title         string       `db:"title"`
	Link          string       `db:"link"`
	Summary       string       `db:"summary"`
	PublishedDate time.Time    `db:"published_at"`
	PostedDate    sql.NullTime `db:"posted_at"`
	CreatedDate   time.Time    `db:"created_at"`
}

type ArticleRepository struct {
	db *pgxpool.Pool
}

func NewArticleRepository(db *pgxpool.Pool) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (a *ArticleRepository) Store(ctx context.Context, article models.Article) error {
	const (
		query = `
			INSERT INTO articles (source_id, title, link, summary, published_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT DO NOTHING;`
	)

	_, err := a.db.Exec(ctx, query, article.SourceID, article.Title, article.Link, article.Summary, article.PublishedDate)
	if err != nil {
		return fmt.Errorf("insert article: %w", err)
	}

	return nil
}

func (a *ArticleRepository) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]*models.Article, error) {
	const (
		query = `
			SELECT id, source_id, title, link, summary, published_at, created_at, posted_at
			FROM articles 
			WHERE posted_at IS NULL AND published_at >= $1::TIMESTAMP
			ORDER BY published_at DESC
			LIMIT $2;`
	)

	rows, err := a.db.Query(ctx, query, since.UTC().Format(time.RFC3339), limit)
	if err != nil {
		return nil, fmt.Errorf("select articles: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		var article dbArticle
		if err := rows.Scan(
			&article.ID,
			&article.SourceID,
			&article.Title,
			&article.Link,
			&article.Summary,
			&article.PublishedDate,
			&article.CreatedDate,
			&article.PostedDate); err != nil {
			return nil, err
		}

		articles = append(articles, &models.Article{
			ID:            article.ID,
			SourceID:      article.SourceID,
			Title:         article.Title,
			Link:          article.Link,
			Summary:       article.Summary,
			PublishedDate: article.PublishedDate,
			PostedDate:    article.PostedDate.Time,
			CreatedDate:   article.CreatedDate,
		})
	}

	return articles, nil
}

func (a *ArticleRepository) MarkPosted(ctx context.Context, id int64) error {
	const (
		query = `UPDATE articles SET posted_at = NOW() WHERE id = $1;`
	)

	_, err := a.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("update article: %w", err)
	}

	return nil
}
