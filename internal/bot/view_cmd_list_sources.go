package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/to77e/news-bot/internal/botkit"
	"github.com/to77e/news-bot/internal/models"
	"strings"
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
