// Package models provides models for the application.
package models

import "time"

// Item - item model.
type Item struct {
	Title      string
	Categories []string
	Link       string
	Date       time.Time
	Summary    string
	SourceName string
}

// Source - source model.
type Source struct {
	ID          int64
	Name        string
	FeedURL     string
	CreatedDate time.Time
}

// Article - article.go model.
type Article struct {
	ID            int64
	SourceID      int64
	Title         string
	Link          string
	Summary       string
	PublishedDate time.Time
	PostedDate    time.Time
	CreatedDate   time.Time
}
