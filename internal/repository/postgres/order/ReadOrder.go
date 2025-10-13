package orderpostgresql

import (
	"context"
	"fmt"

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
		SELECT * from orders WHERE order_uuid = $1
	`, orderID).Scan(&order.OrderUUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		fmt.Println(err)
		return nil, types.ErrOrderNotFound
	}

	err = tx.QueryRowContext(ctx, `
		SELECT name, phone, zip, city, address, region, email
		FROM delivery WHERE order_uuid = $1
	`, orderID).Scan(&delivery.Name, &delivery.Phone, &delivery.Zip,
		&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil {
		return nil, types.ErrDeliveryNotFound
	}

	err = tx.QueryRowContext(ctx, `
		SELECT transaction, request_id, currency, provider, amount, 
		payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment WHERE order_uuid = $1
	`, orderID).Scan(&payment.Transaction, &payment.RequestID, &payment.Currency,
		&payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank,
		&payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)
	if err != nil {
		return nil, types.ErrPaymentNotFound
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT chrt_id, track_number, price, rid, name, sale, size, 
		       total_price, nm_id, brand, status
		FROM items 
		WHERE order_uuid = $1
	`, orderID)
	if err != nil {
		return nil, types.ErrItemsNotFound
	}
	defer rows.Close()

	for rows.Next() {
		var item types.Item
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price,
			&item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
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
