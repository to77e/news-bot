package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/to77e/news-bot/internal/botkit"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(
			update.FromChat().ID,
			"Hello, I'm a news bot. Type /help to see the list of commands.")); err != nil {
			return fmt.Errorf("send message: %w", err)
		}
		return nil
	}
}
