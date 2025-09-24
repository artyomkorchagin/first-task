package types

import "context"

type ServiceInterface interface {
	ReadOrder(ctx context.Context, id string) error
	CreateOrder(ctx context.Context, order *Order) error
}
