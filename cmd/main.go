package main

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/to77e/news-bot/internal/bot"
	"github.com/to77e/news-bot/internal/botkit"
	"github.com/to77e/news-bot/internal/config"
	"github.com/to77e/news-bot/internal/database"
	"github.com/to77e/news-bot/internal/fetcher"
	"github.com/to77e/news-bot/internal/notifier"
	"github.com/to77e/news-bot/internal/repository"
	"github.com/to77e/news-bot/internal/summary"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const logLevel = slog.LevelInfo

	ctx := context.Background()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	if err := config.Read(".config.yaml"); err != nil {
		slog.With("error", err.Error()).ErrorContext(ctx, "read config")
		return
	}
	cfg := config.Get()

	botAPI, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		slog.With("error", err.Error()).ErrorContext(ctx, "create bot")
		return
	}

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)
	var conn *pgxpool.Pool
	conn, err = database.NewPostgres(ctx, dsn)
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
