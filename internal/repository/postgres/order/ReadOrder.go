package orderpostgresql

import (
	"context"
	"fmt"

	"github.com/artyomkorchagin/first-task/internal/types"
)

func (r *Repository) ReadOrder(ctx context.Context, orderID string) (*types.Order, error) {
	return nil, fmt.Errorf("not implemented")
}
