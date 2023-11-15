package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/to77e/news-fetching-bot/internal/botkit"
	"strings"
)

func ViewCmdInfo(version, commitHash string) botkit.ViewFunc {
	message := strings.Builder{}
	message.WriteString("Version: ")
	message.WriteString(version)
	message.WriteString("\n")
	message.WriteString("Commit Hash: ")
	message.WriteString(commitHash)

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(
			update.FromChat().ID,
			message.String())); err != nil {
			return fmt.Errorf("send message: %w", err)
		}
		return nil
	}
}
