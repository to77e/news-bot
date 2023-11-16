package bot

import (
	"context"
	"fmt"
	"github.com/to77e/news-fetching-bot/internal/botkit/markup"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/to77e/news-fetching-bot/internal/botkit"
	"github.com/to77e/news-fetching-bot/internal/models"
)

type SourceLister interface {
	Sources(ctx context.Context) ([]*models.Source, error)
}

func ViewCmdListSources(lister SourceLister) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		sources, err := lister.Sources(ctx)
		if err != nil {
			return fmt.Errorf("list sources: %w", err)
		}

		var sourceInfos []string
		for _, v := range sources {
			src := formatSource(v)
			sourceInfos = append(sourceInfos, src)
		}
		msgText := fmt.Sprintf(
			"List sources \\(total %d\\):\n\n%s",
			len(sources),
			strings.Join(sourceInfos, "\n\n"),
		)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		reply.ParseMode = parseModeMarkdownV2

		if _, err := bot.Send(reply); err != nil {
			return fmt.Errorf("send message: %w", err)
		}

		return nil
	}
}

func formatSource(source *models.Source) string {
	return fmt.Sprintf(
		"*%s*\nID: `%d`\nfeed URL: %s",
		markup.EscapeForMarkdown(source.Name),
		source.ID,
		markup.EscapeForMarkdown(source.URL),
	)
}
