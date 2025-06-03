package models

import (
	"github.com/google/uuid"
	"time"
)

type UserSubscription struct {
	ID                 uint           `gorm:"primaryKey"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null"` // Изменяем на uuid
	PlanID             string         `gorm:"type:varchar(255);not null"`
	StripeSubscriptionID string        `gorm:"type:varchar(255);not null;unique"`
	Status             string         `gorm:"type:varchar(50);not null"`
	CurrentPeriodStart time.Time      `gorm:"not null"`
	CurrentPeriodEnd   time.Time      `gorm:"not null"`
	CreatedAt          time.Time      `gorm:"autoCreateTime"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime"`
}

func (UserSubscription) TableName() string {
	return "user_subscriptions"
}