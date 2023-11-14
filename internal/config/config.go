package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

const (
	version    = "dev"
	commitHash = "-"
)

var cfg *Config

type Config struct {
	Project  Project  `yaml:"project"`
	Settings Settings `yaml:"settings"`
	Telegram Telegram `yaml:"telegram"`
	Database Database `yaml:"database"`
	OpenAI   OpenAI   `yaml:"openai"`
}

type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string
	CommitHash  string
}

type Settings struct {
	FetchInterval        time.Duration `yaml:"fetch_interval"`
	NotificationInterval time.Duration `yaml:"notification_interval"`
	FilterKeyword        []string      `yaml:"filter_keyword"`
}

type Telegram struct {
	BotToken  string `yaml:"bot_token"`
	ChannelID int64  `yaml:"channel_id"`
}

type Database struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Migrations string `yaml:"migrations"`
	Name       string `yaml:"name"`
	SslMode    string `yaml:"sslmode"`
	Driver     string `yaml:"driver"`
}

type OpenAI struct {
	Key    string `yaml:"key"`
	Prompt string `yaml:"prompt"`
}

func Get() Config {
	if cfg != nil {
		return *cfg
	}
	return Config{}
}

func Read(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return fmt.Errorf("open config file: %w", err)
	}
	defer file.Close() //nolint: errcheck

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
