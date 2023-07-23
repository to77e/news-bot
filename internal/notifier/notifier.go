// Package notifier provides work with notifier.
package notifier

import (
	"context"
	"fmt"
	"github.com/go-shiori/go-readability"
	"github.com/to77e/news-bot/internal/botkit/markup"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/to77e/news-bot/internal/models"
)

// ArticleProvider is a service that provides articles.
type ArticleProvider interface {
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]*models.Article, error)
	MarkPosted(ctx context.Context, id int64) error
}

// Summarizer is a service that summarizes text.
type Summarizer interface {
	Summarize(ctx context.Context, text string) (string, error)
}

// Notifier is a service that sends articles to the channel.
type Notifier struct {
	articles         ArticleProvider
	summarizer       Summarizer
	bot              *tgbotapi.BotAPI
	sendInterval     time.Duration
	lookupTimeWindow time.Duration
	channelID        int64
}

// New creates a new notifier.
func New(
	articles ArticleProvider,
	summarizer Summarizer,
	bot *tgbotapi.BotAPI,
	sendInterval time.Duration,
	lookupTimeWindow time.Duration,
	channelID int64,
) *Notifier {
	return &Notifier{
		articles:         articles,
		summarizer:       summarizer,
		bot:              bot,
		sendInterval:     sendInterval,
		lookupTimeWindow: lookupTimeWindow,
		channelID:        channelID,
	}
}

// Start starts the notifier.
func (n *Notifier) Start(ctx context.Context) error {
	ticker := time.NewTicker(n.sendInterval)
	defer ticker.Stop()

	if err := n.SelectAndSendArticle(ctx); err != nil {
		return fmt.Errorf("select and send article: %w", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := n.SelectAndSendArticle(ctx); err != nil {
				return fmt.Errorf("select and send article: %w", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// SelectAndSendArticle selects one article and sends it to the channel.
func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	//TODO: wrap in a transaction
	topOneArticles, err := n.articles.AllNotPosted(ctx, time.Now().Add(-n.lookupTimeWindow), 1)
	if err != nil {
		return fmt.Errorf("failed to get top one article: %w", err)
	}

	if len(topOneArticles) == 0 {
		return nil
	}

	article := topOneArticles[0]
	summary, err := n.extractSummary(ctx, article)
	if err != nil {
		return fmt.Errorf("failed to extract summary: %w", err)
	}

	if err := n.sendArticle(article, summary); err != nil {
		return fmt.Errorf("failed to send article: %w", err)
	}

	return n.articles.MarkPosted(ctx, article.ID)
}

func (n *Notifier) extractSummary(ctx context.Context, article *models.Article) (string, error) {
	var r io.Reader

	if article.Summary != "" {
		r = strings.NewReader(article.Summary)
	} else {
		resp, err := http.Get(article.Link)
		if err != nil {
			return "", fmt.Errorf("failed to get article: %w", err)
		}
		defer resp.Body.Close()

		r = resp.Body
	}

	doc, err := readability.FromReader(r, nil)
	if err != nil {
		return "", fmt.Errorf("failed to parse article: %w", err)
	}

	summary, err := n.summarizer.Summarize(ctx, cleanText(doc.TextContent))
	if err != nil {
		return "", fmt.Errorf("failed to summarize: %w", err)
	}

	return fmt.Sprintf("\n\n%s", summary), nil
}

func (n *Notifier) sendArticle(article *models.Article, summary string) error {
	const (
		messageFormat = "*%s*%s\n\n%s"
	)

	msg := tgbotapi.NewMessage(n.channelID, fmt.Sprintf(
		messageFormat,
		markup.EscapeMarkdown(article.Title),
		markup.EscapeMarkdown(summary),
		markup.EscapeMarkdown(article.Link)),
	)
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := n.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

var redundantNewLines = regexp.MustCompile(`\n{3,}`)

func cleanText(text string) string {
	return redundantNewLines.ReplaceAllString(text, "\n")
}
