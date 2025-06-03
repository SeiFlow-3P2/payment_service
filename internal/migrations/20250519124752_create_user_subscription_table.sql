-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    plan_id INTEGER NOT NULL,
    stripe_subscription_id TEXT,
    status TEXT CHECK (status IN ('active', 'inactive', 'canceled')) NOT NULL,
    current_period_start TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    current_period_end TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_subscription;
-- +goose StatementEnd
