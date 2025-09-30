package orderpostgresql

import (
	"context"

	"github.com/artyomkorchagin/first-task/internal/types"
)

func (r *Repository) CreateOrder(ctx context.Context, o *types.Order) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO orders VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, o.OrderUUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature,
		o.CustomerID, o.DeliveryService, o.Shardkey, o.SmID, o.DateCreated, o.OofShard)
	if err != nil {
		return types.ErrDB
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO delivery VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, o.OrderUUID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip,
		o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email)
	if err != nil {
		return types.ErrDB
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO payment VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, o.OrderUUID, o.Payment.Transaction, o.Payment.RequestID, o.Payment.Currency,
		o.Payment.Provider, o.Payment.Amount, o.Payment.PaymentDt, o.Payment.Bank,
		o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee)
	if err != nil {
		return types.ErrDB
	}

	for _, item := range o.Items {
		if _, err = r.db.ExecContext(ctx, `
			INSERT INTO items VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, o.OrderUUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmID, item.Brand, item.Status); err != nil {
			return types.ErrDB
		}
	}
	tx.Commit()
	return nil
}
