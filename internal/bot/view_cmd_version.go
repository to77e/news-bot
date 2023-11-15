package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/to77e/news-fetching-bot/internal/botkit"
)

func ViewCmdVersion(version string) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(
			update.FromChat().ID,
			version)); err != nil {
			return fmt.Errorf("send message: %w", err)
		}
		return nil
	}
}
