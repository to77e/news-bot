package summary

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
)

type OpenAISummarizer struct {
	client  *openai.Client
	prompt  string
	model   string
	enabled bool
	mu      sync.Mutex
}

func NewOpenAISummarizer(apiKey, prompt, model string) *OpenAISummarizer {
	s := &OpenAISummarizer{
		client: openai.NewClient(apiKey),
		prompt: prompt,
		model:  model,
	}

	slog.Info("openai summarizer", "is enabled", apiKey != "")

	if apiKey != "" {
		s.enabled = true
	}

	return s
}

func (s *OpenAISummarizer) Summarize(ctx context.Context, text string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.enabled {
		return "", nil
	}

	request := openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: s.prompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
		MaxTokens:   1024,
		Temperature: 1,
		TopP:        1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "status code: 429") {
			slog.Warn("openai summarizer", "rate limit exceeded", err)
			return "", nil
		}
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in openai response")
	}

	raw := strings.TrimSpace(resp.Choices[0].Message.Content)
	if strings.HasSuffix(raw, ".") {
		return raw, nil
	}

	sentences := strings.Split(raw, ".")

	return strings.Join(sentences[:len(sentences)-1], ".") + ".", nil
}
