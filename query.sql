-- name: CreateOrder :exec
INSERT INTO orders (
    order_uid, track_number, entry, locale, internal_signature,
    customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
);

-- name: CreateDelivery :exec
INSERT INTO delivery (
    order_uid, name, phone, zip, city, address, region, email
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: CreatePayment :exec
INSERT INTO payment (
    order_uid, transaction, request_id, currency, provider, amount,
    payment_dt, bank, delivery_cost, goods_total, custom_fee
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
);

-- name: CreateItem :exec
INSERT INTO items (
    order_uid, chrt_id, track_number, price, rid, name, sale, size,
    total_price, nm_id, brand, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);

-- name: GetOrder :one
SELECT * FROM orders WHERE order_uid = $1;

-- name: GetDelivery :one
SELECT * FROM delivery WHERE order_uid = $1;

-- name: GetPayment :one
SELECT * FROM payment WHERE order_uid = $1;

-- name: GetItems :many
SELECT * FROM items WHERE order_uid = $1;