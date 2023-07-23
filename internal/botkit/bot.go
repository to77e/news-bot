// Package botkit provides work with bot.
package botkit

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"runtime/debug"
	"time"
)

// Bot is a bot.
type Bot struct {
	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
}

// ViewFunc is a view function.
type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

// New creates a new bot.
func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		api: api,
	}
}

// RegisterCmdView registers a view function for the command.
func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}
	b.cmdViews[cmd] = view
}

// Run runs the bot.
func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateContext, updateChannel := context.WithTimeout(ctx, 5*time.Minute)
			b.handleUpdate(updateContext, update)
			updateChannel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic recovery: %v\n%s", p, string(debug.Stack()))
		}
	}()

	if (update.Message == nil || !update.Message.IsCommand()) && update.CallbackQuery == nil {
		return
	}
	var view ViewFunc
	if !update.Message.IsCommand() {
		return
	}
	cmd := update.Message.Command()
	cmdView, ok := b.cmdViews[cmd]
	if !ok {
		return
	}
	view = cmdView
	if err := view(ctx, b.api, update); err != nil {
		log.Printf("handling update: %v", err)

		if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "internal error")); err != nil {
			log.Printf("send message: %v", err)
		}
	}
}
