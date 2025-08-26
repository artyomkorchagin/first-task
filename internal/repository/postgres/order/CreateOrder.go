package orderpostgresql

import (
	"context"
	"fmt"

	"github.com/artyomkorchagin/first-task/internal/types"
)

func (r *Repository) CreateOrder(ctx context.Context, o *types.Order) error {
	return fmt.Errorf("not implemented")
}
