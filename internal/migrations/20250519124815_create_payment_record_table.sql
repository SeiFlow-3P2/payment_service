-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payment_record (
    id UUID PRIMARY KEY,
    user_id INTEGER NOT NULL,
    user_subscription_id INTEGER,
    stripe_charge_id VARCHAR(255),
    stripe_checkout_session_id VARCHAR(255),
    amount_id DECIMAL(10, 2),
    currency_id VARCHAR(10),
    status_id VARCHAR(50),
    payment_method_details_id JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_record;
-- +goose StatementEnd
