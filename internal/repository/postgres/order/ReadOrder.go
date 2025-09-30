package orderpostgresql

import (
	"context"

	"github.com/artyomkorchagin/first-task/internal/types"
)

func (r *Repository) ReadOrder(ctx context.Context, orderID string) (*types.Order, error) {
	var (
		order    types.Order
		delivery types.Delivery
		payment  types.Payment
		items    []types.Item
	)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
		SELECT from orders WHERE order_uid = $1
	`, orderID).Scan(&order)
	if err != nil {
		return nil, types.ErrOrderNotFound
	}

	err = tx.QueryRowContext(ctx, `
		SELECT from delivery WHERE order_uid = $1
	`, orderID).Scan(&delivery)
	if err != nil {
		return nil, types.ErrDeliveryNotFound
	}

	err = tx.QueryRowContext(ctx, `
		SELECT from payment WHERE order_uid = $1
	`, orderID).Scan(&payment)
	if err != nil {
		return nil, types.ErrPaymentNotFound
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT from items WHERE order_uid = $1
	`, orderID)
	if err != nil {
		return nil, types.ErrItemsNotFound
	}
	for rows.Next() {
		var item types.Item
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	tx.Commit()
	order.Delivery = delivery
	order.Payment = payment
	order.Items = items
	return &order, nil
}
