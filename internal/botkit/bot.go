package botkit

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		api: api,
	}
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}
	b.cmdViews[cmd] = view
}

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
			slog.With("panic", p, "stacktrace", string(debug.Stack())).Error("panic recovery")
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
		slog.With("error", err.Error()).ErrorContext(ctx, "handling update")

		if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "internal error")); err != nil {
			slog.With("error", err.Error()).ErrorContext(ctx, "send message")
		}
	}
}

func (b *Bot) GetCommandNames() []string {
	var names []string
	for name := range b.cmdViews {
		names = append(names, name)
	}
	return names
}
