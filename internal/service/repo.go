package orderservice

import (
	"context"

	"github.com/artyomkorchagin/first-task/internal/types"
)

type ReadWriter interface {
	Reader
	Writer
}

type Reader interface {
	ReadOrder(ctx context.Context, id string) (*types.Order, error)
}

type Writer interface {
	CreateOrder(ctx context.Context, order *types.Order) error
}
