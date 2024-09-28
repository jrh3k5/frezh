package chatgpt

import "context"

type Service interface {
	Ask(ctx context.Context, question string) (string, error)
}
