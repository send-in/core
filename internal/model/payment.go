package model

import (
	"time"
	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSucceeded PaymentStatus = "succeeded"
	PaymentFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID          uuid.UUID `gorm:"<-:create;type:uuid;default:gen_random_uuid();primaryKey"`
	Status      PaymentStatus `gorm:"default:'pending'"`
	AccountID   *uuid.UUID `gorm:"not null"`
	
	PlanCredits int
	Amount      int64
	Provider    string
	ExternalID  string
	Account     Account
	CreatedAt   time.Time
	UpdatedAt   time.Time
}