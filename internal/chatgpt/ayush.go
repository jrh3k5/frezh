package chatgpt

import (
	"context"
	"fmt"

	"github.com/ayush6624/go-chatgpt"
)

type AyushService struct {
	client *chatgpt.Client
}

func NewAyushService(client *chatgpt.Client) *AyushService {
	return &AyushService{client: client}
}

func (s *AyushService) Ask(ctx context.Context, question string) (string, error) {
	response, err := s.client.Send(ctx, &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT35Turbo,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: question,
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to send message to ChatGPT: %w", err)
	}

	return response.Choices[0].Message.Content, nil
}
