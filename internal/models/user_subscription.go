package models

import "time"

type UserSubscription struct {
	ID                   uint      `gorm:"primaryKey;column:id"`
	UserID               uint      `gorm:"column:user_id"`
	PlanID               string    `gorm:"column:plan_id"`
	StripeSubscriptionID string    `gorm:"column:stripe_subscription_id"`
	Status               string    `gorm:"column:status"`
	CurrentPeriodStart   time.Time `gorm:"column:current_period_start"`
	CurrentPeriodEnd     time.Time `gorm:"column:current_period_end"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}
