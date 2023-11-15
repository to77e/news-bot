package fetcher

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/to77e/news-fetching-bot/internal/models"
	"github.com/to77e/news-fetching-bot/internal/source"
)

type ArticleRepository interface {
	Store(ctx context.Context, article models.Article) error
}

type SourceRepository interface {
	Sources(ctx context.Context) ([]*models.Source, error)
}

type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]models.Item, error)
}

type Fetcher struct {
	articles ArticleRepository
	sources  SourceRepository

	fetchInterval time.Duration
	filterKeyword []string
}

func New(
	articles ArticleRepository,
	sources SourceRepository,
	fetchInterval time.Duration,
	filterKeyword []string,
) *Fetcher {
	return &Fetcher{
		articles:      articles,
		sources:       sources,
		fetchInterval: fetchInterval,
		filterKeyword: filterKeyword,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return fmt.Errorf("fetch: %w", err)
			}
		}
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sources, err := f.sources.Sources(ctx)
	if err != nil {
		return fmt.Errorf("fetch sources: %w", err)
	}

	var wg sync.WaitGroup
	for _, v := range sources {
		wg.Add(1)
		rssSource := source.NewRSSSourceForModel(v)

		go func(source Source) {
			defer wg.Done()
			items, err := source.Fetch(ctx)
			if err != nil {
				slog.With("error", err.Error()).ErrorContext(ctx, "fetch source", "name", source.Name())
				return
			}

			if err := f.processItems(ctx, source, items); err != nil {
				slog.With("error", err.Error()).ErrorContext(ctx, "process items", "name", source.Name())
				return
			}

		}(rssSource)
	}

	wg.Wait()
	return nil
}

func (f *Fetcher) processItems(ctx context.Context, source Source, items []models.Item) error {
	for _, v := range items {
		v.Date = v.Date.UTC()

		if f.itemShouldBeSkipped(v) {
			continue
		}
		if err := f.articles.Store(ctx, models.Article{
			SourceID:      source.ID(),
			Title:         v.Title,
			Link:          v.Link,
			Summary:       v.Summary,
			PublishedDate: v.Date,
		}); err != nil {
			return fmt.Errorf("store article.go: %w", err)
		}
	}
	return nil
}

func (f *Fetcher) itemShouldBeSkipped(item models.Item) bool {
	var categoryContainsKeyword bool
	for _, keyword := range f.filterKeyword {
		for _, category := range item.Categories {
			if strings.Contains(strings.ToLower(category), keyword) {
				categoryContainsKeyword = true
				break
			}
		}
		titleContainsKeyword := strings.Contains(strings.ToLower(item.Title), keyword)
		if categoryContainsKeyword || titleContainsKeyword {
			return true
		}
	}
	return false
}
