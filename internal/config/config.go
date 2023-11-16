package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"time"
)

var (
	version    = "dev"
	commitHash = "-"
)

var cfg *Config

type Config struct {
	Project  Project
	Settings Settings
	Telegram Telegram
	Database Database
	OpenAI   OpenAI
}

type Project struct {
	Name        string `env:"PROJECT_NAME"`
	Environment string `env:"PROJECT_ENVIRONMENT" envDefault:"development"`
	LogLevel    int    `env:"PROJECT_LOG_LEVEL" envDefault:"0"`
	Version     string
	CommitHash  string
}

type Settings struct {
	FetchInterval        time.Duration `env:"FETCH_INTERVAL"`
	NotificationInterval time.Duration `env:"NOTIFICATION_INTERVAL"`
	FilterKeyword        []string
}

type Telegram struct {
	BotToken  string `env:"TELEGRAM_BOT_TOKEN"`
	ChannelID int64  `env:"TELEGRAM_CHANNEL_ID"`
}

type Database struct {
	Host     string `env:"DATABASE_HOST"`
	Port     string `env:"DATABASE_PORT"`
	User     string `env:"DATABASE_USER"`
	Password string `env:"DATABASE_PASSWORD"`
	Name     string `env:"DATABASE_NAME"`
	SSLMode  string `env:"DATABASE_SSL_MODE" envDefault:"disable"`
}

type OpenAI struct {
	Key    string `env:"OPENAI_API_KEY"`
	Prompt string `env:"OPENAI_API_PROMPT"`
	Model  string `env:"OPENAI_API_MODEL" envDefault:"gpt-3.5-turbo"`
}

func Get() Config {
	if cfg != nil {
		return *cfg
	}
	return Config{}
}

func Read() error {
	if cfg != nil {
		return nil
	}

	cfg = &Config{}
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
