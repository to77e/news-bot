package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/to77e/news-fetching-bot/internal/bot"
	"github.com/to77e/news-fetching-bot/internal/botkit"
	"github.com/to77e/news-fetching-bot/internal/config"
	"github.com/to77e/news-fetching-bot/internal/database"
	"github.com/to77e/news-fetching-bot/internal/fetcher"
	"github.com/to77e/news-fetching-bot/internal/notifier"
	"github.com/to77e/news-fetching-bot/internal/repository"
	"github.com/to77e/news-fetching-bot/internal/summary"
)

func main() {
	ctx := context.Background()

	if err := config.Read(); err != nil {
		slog.With("error", err.Error()).ErrorContext(ctx, "read config")
		return
	}
	cfg := config.Get()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(cfg.Project.LogLevel)})))

	botAPI, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		slog.With("error", err.Error()).ErrorContext(ctx, "create bot")
		return
	}

	var conn *pgxpool.Pool
	conn, err = database.NewPostgres(
		ctx,
		fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		),
	)
	if err != nil {
		slog.With("error", err.Error()).ErrorContext(ctx, "new postgresql connection")
		return
	}
	defer conn.Close()

	articleRepository := repository.NewArticleRepository(conn)
	sourceRepository := repository.NewSourceRepository(conn)
	fetch := fetcher.New(
		articleRepository,
		sourceRepository,
		cfg.Settings.FetchInterval,
		cfg.Settings.FilterKeyword,
	)
	notify := notifier.New(
		articleRepository,
		summary.NewOpenAISummarizer(cfg.OpenAI.Key, cfg.OpenAI.Prompt),
		botAPI,
		cfg.Settings.NotificationInterval,
		2*cfg.Settings.FetchInterval,
		cfg.Telegram.ChannelID,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	newsBot := botkit.New(botAPI)
	newsBot.RegisterCmdView("start", bot.ViewCmdStart())
	newsBot.RegisterCmdView("add_source", bot.ViewCmdAddSource(sourceRepository))
	newsBot.RegisterCmdView("list_sources", bot.ViewCmdListSources(sourceRepository))
	newsBot.RegisterCmdView("help", bot.ViewCmdHelp(newsBot.GetCommandNames()))
	newsBot.RegisterCmdView("version", bot.ViewCmdVersion(cfg.Project.Version))

	go func(ctx context.Context) {
		if err := fetch.Start(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				slog.With("error", err.Error()).Error("fetcher start")
				return
			}
			slog.With("error", err.Error()).Error("fetcher stop")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := notify.Start(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				slog.With("error", err.Error()).Error("notifier start")
				return
			}
			slog.With("error", err.Error()).Error("notifier stop")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			slog.With("error", err.Error()).Error("bot start")
			return
		}
		slog.With("error", err.Error()).Error("bot stop")
	}
}
