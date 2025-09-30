-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_uuid TEXT PRIMARY KEY,
    track_number TEXT NOT NULL,
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INTEGER,
    date_created TIMESTAMPTZ,
    oof_shard TEXT
);

CREATE TABLE IF NOT EXISTS delivery (
    order_uuid TEXT PRIMARY KEY REFERENCES orders(order_uuid) ON DELETE CASCADE,
    name TEXT NOT NULL,
    phone TEXT,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT
);

CREATE TABLE IF NOT EXISTS payment (
    order_uuid TEXT PRIMARY KEY REFERENCES orders(order_uuid) ON DELETE CASCADE,
    transaction TEXT NOT NULL,
    request_id TEXT,
    currency TEXT,
    provider TEXT,
    amount INTEGER,
    payment_dt INTEGER, 
    bank TEXT,
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);

CREATE TABLE IF NOT EXISTS items (
    order_uuid TEXT PRIMARY KEY REFERENCES orders(order_uuid) ON DELETE CASCADE,
    chrt_id BIGINT,
    track_number TEXT,
    price INTEGER,
    rid TEXT,
    name TEXT,
    sale INTEGER,
    size TEXT,
    total_price INTEGER,
    nm_id INTEGER,
    brand TEXT,
    status INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;

DROP TABLE IF EXISTS delivery;

DROP TABLE IF EXISTS payment;

DROP TABLE IF EXISTS items;
-- +goose StatementEnd
