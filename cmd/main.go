// Package main - entry point for the application.
package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/to77e/news-bot/internal/bot"
	"github.com/to77e/news-bot/internal/botkit"
	"github.com/to77e/news-bot/internal/config"
	"github.com/to77e/news-bot/internal/database"
	"github.com/to77e/news-bot/internal/fetcher"
	"github.com/to77e/news-bot/internal/notifier"
	"github.com/to77e/news-bot/internal/repository"
	"github.com/to77e/news-bot/internal/summary"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("create bot: %v", err)
		return
	}

	conn, err := database.NewPostgres(ctx, config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("create connection: %v", err)
		return
	}
	defer conn.Close()

	articleRepository := repository.NewArticleRepository(conn)
	sourceRepository := repository.NewSourceRepository(conn)
	fetch := fetcher.New(
		articleRepository,
		sourceRepository,
		config.Get().FetchInterval,
		config.Get().FilterKeyword,
	)
	notify := notifier.New(
		articleRepository,
		summary.NewOpenAISummarizer(config.Get().OpenAIKey, config.Get().OpenAIPrompt),
		botAPI,
		config.Get().NotificationInterval,
		2*config.Get().FetchInterval,
		config.Get().TelegramChannelID,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	newsBot := botkit.New(botAPI)
	newsBot.RegisterCmdView("start", bot.ViewCmdStart())
	newsBot.RegisterCmdView("add_source", bot.ViewCmdAddSource(sourceRepository))
	newsBot.RegisterCmdView("list_sources", bot.ViewCmdListSources(sourceRepository))

	go func(ctx context.Context) {
		if err := fetch.Start(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] fetcher start: %v", err)
				return
			}
			log.Printf("fetcher stop: %v", err)
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := notify.Start(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] notifier start: %v", err)
				return
			}
			log.Printf("notifier stop: %v", err)
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] bot start: %v", err)
			return
		}
		log.Printf("bot stop: %v", err)
	}
}
