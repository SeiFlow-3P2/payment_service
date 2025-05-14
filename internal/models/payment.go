package models

import (
	"time"
)

type Payment struct {
	IDPayment             uint      `gorm:"primaryKey;column:id_payment"`
	UserID                string    `gorm:"column:user_id"`
	StripeID              string    `gorm:"column:stripe_id"`
	StripeSubscriptionID  string    `gorm:"column:stripe_subscription_id"`
	Status                string    `gorm:"column:status"`
	PriceID               string    `gorm:"column:price_id"`
	Currency              string    `gorm:"column:currency"`
	CurrentPeriodEnd      time.Time `gorm:"column:current_period_end"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
