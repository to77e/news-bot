package config

import (
	"github.com/cristalhq/aconfig/aconfigyaml"
	"log"
	"sync"
	"time"

	"github.com/cristalhq/aconfig"
)

type Config struct {
	TelegramBotToken     string        `yaml:"telegram_bot_token" env:"TELEGRAM_BOT_TOKEN" required:"true"`
	TelegramChannelID    int64         `yaml:"telegram_channel_id" env:"TELEGRAM_CHANNEL_ID" required:"true"`
	DatabaseDSN          string        `yaml:"database_dsn" env:"DATABASE_DSN" default:"postgres://postgres:postgres@localhost:5432/news?sslmode=disable"`
	FetchInterval        time.Duration `yaml:"fetch_interval" env:"FETCH_INTERVAL" default:"10m"`
	NotificationInterval time.Duration `yaml:"notification_interval" env:"NOTIFICATION_INTERVAL" default:"1h"`
	FilterKeyword        []string      `yaml:"filter_keyword" env:"FILTER_KEYWORD"`
	OpenAIKey            string        `yaml:"openai_key" env:"OPENAI_KEY"`
	OpenAIPrompt         string        `yaml:"openai_prompt" env:"OPENAI_PROMPT"`
}

var (
	cfg  *Config
	once sync.Once
)

func Get() *Config {
	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "NB",
			Files:     []string{"./config.yaml", "./config.local.yaml"},
			FileDecoders: map[string]aconfig.FileDecoder{
				".yaml": aconfigyaml.New(),
			},
		})
		if err := loader.Load(); err != nil {
			log.Printf("[ERROR] load config: %v", err)
		}
	})
	return cfg
}
