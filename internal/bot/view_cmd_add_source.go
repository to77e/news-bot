package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/to77e/news-fetching-bot/internal/botkit"
	"github.com/to77e/news-fetching-bot/internal/models"
)

const parseModeMarkdownV2 = "MarkdownV2"

type SourceRepository interface {
	Add(ctx context.Context, source models.Source) (int64, error)
}

func ViewCmdAddSource(storage SourceRepository) botkit.ViewFunc {
	type addSourceArgs struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Priority int    `json:"priority"`
	}

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[addSourceArgs](update.Message.CommandArguments())
		if err != nil {
			return fmt.Errorf("parse JSON: %w", err)
		}

		source := models.Source{
			Name: args.Name,
			URL:  args.URL,
		}

		sourceID, err := storage.Add(ctx, source)
		if err != nil {
			return fmt.Errorf("add source: %w", err)
		}

		var (
			msgText = fmt.Sprintf(
				"Source added with ID: `%d`\\. Use this ID for updating the source or deleting it\\.",
				sourceID,
			)
			reply = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)
		reply.ParseMode = parseModeMarkdownV2

		if _, err := bot.Send(reply); err != nil {
			return fmt.Errorf("send message: %w", err)
		}

		return nil
	}
}
