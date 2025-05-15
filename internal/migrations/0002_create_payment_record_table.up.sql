CREATE TABLE IF NOT EXISTS payment_record (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    user_subscription_id INTEGER,
    stripe_charge_id VARCHAR(255),
    stripe_checkout_session_id VARCHAR(255),
    amount DECIMAL(10, 2),
    currency VARCHAR(10),
    status VARCHAR(50),
    payment_method_details JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

