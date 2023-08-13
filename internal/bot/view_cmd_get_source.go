package bot

import (
	"fmt"
	"github.com/to77e/news-bot/internal/botkit/markup"
	"github.com/to77e/news-bot/internal/models"
)

func formatSource(source *models.Source) string {
	return fmt.Sprintf(
		"*%s*\nID: `%d`\nfeed URL: %s",
		markup.EscapeForMarkdown(source.Name),
		source.ID,
		markup.EscapeForMarkdown(source.FeedURL),
	)
}
