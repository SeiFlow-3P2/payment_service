package models

import "time"

type PaymentRecord struct {
	ID                        uint      `gorm:"primaryKey;column:id"`
	UserID                    string    `gorm:"column:user_id"`
	UserSubscriptionID        uint      `gorm:"column:user_subscription_id"`
	StripeChargeID            string    `gorm:"column:stripe_charge_id"`
	StripeCheckoutSessionID   string    `gorm:"column:stripe_checkout_session_id"`
	Amount                    int64     `gorm:"column:amount"`
	Currency                  string    `gorm:"column:currency"`
	Status                    string    `gorm:"column:status"`
	PaymentMethodDetails      string    `gorm:"type:jsonb;column:payment_method_details"`
	CreatedAt                 time.Time `gorm:"column:created_at"`
	UpdatedAt                 time.Time `gorm:"column:updated_at"`
}
