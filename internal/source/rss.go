package source

import (
	"context"
	"fmt"

	"github.com/SlyMarbo/rss"
	"github.com/to77e/news-fetching-bot/internal/models"
)

type RSSSource struct {
	URL        string
	SourceID   int64
	SourceName string
}

func NewRSSSourceForModel(m *models.Source) RSSSource {
	return RSSSource{
		URL:        m.URL,
		SourceID:   m.ID,
		SourceName: m.Name,
	}
}

func (r RSSSource) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := r.loadFeed(ctx, r.URL)
	if err != nil {
		return nil, fmt.Errorf("fetch url %s: %w", r.URL, err)
	}

	var items []models.Item
	for _, v := range feed.Items {
		items = append(items, models.Item{
			Title:      v.Title,
			Categories: v.Categories,
			Link:       v.Link,
			Date:       v.Date,
			Summary:    v.Summary,
			SourceName: r.SourceName,
		})
	}
	return items, nil
}

func (r RSSSource) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	var (
		feedCh = make(chan *rss.Feed)
		errCh  = make(chan error)
	)

	go func() {
		feed, err := rss.Fetch(url)
		if err != nil {
			errCh <- err
			return
		}

		feedCh <- feed
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case feed := <-feedCh:
		return feed, nil
	}
}

func (r RSSSource) ID() int64 {
	return r.SourceID
}

func (r RSSSource) Name() string {
	return r.SourceName
}
