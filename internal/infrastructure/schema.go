package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EnsureSchema(ctx context.Context, db *pgxpool.Pool) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    number TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status TEXT NOT NULL,
    accrual NUMERIC(10,2),
    uploaded_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_orders_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS withdrawals (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    order_number TEXT UNIQUE NOT NULL,
    sum NUMERIC(10,2) NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_withdrawals_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
`
	_, err := db.Exec(ctx, ddl)
	return err
}
