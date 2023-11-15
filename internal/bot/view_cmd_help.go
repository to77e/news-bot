package bot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/to77e/news-fetching-bot/internal/botkit"
)

func ViewCmdHelp(commands []string) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		var message strings.Builder

		message.WriteString("List of commands:\n")
		for _, command := range commands {
			message.WriteString(fmt.Sprintf("/%s\n", command))
		}

		if _, err := bot.Send(tgbotapi.NewMessage(
			update.FromChat().ID,
			message.String())); err != nil {
			return fmt.Errorf("send message: %w", err)
		}
		return nil
	}
}
