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
	ID 			  uuid.UUID `gorm:"<-:create;type:uuid;default:gen_random_uuid();primaryKey"`
	Status		  PaymentStatus `gorm:"default:'pending'"`

	AccountID 	  uuid.UUID `gorm:"not null;index"`
	Account   	  Account
	
	Plan      	  string
	PlanCredits   int
	Amount        int64
	Currency  	  string
	Provider  	  string
	ExternalID	  string `gorm:"uniqueIndex"`
	OrderID   	  string `gorm:"uniqueIndex"`
	Payload       []byte `gorm:"type:jsonb"`

	CompletedAt  *time.Time
	CreatedAt 	  time.Time
	UpdatedAt 	  time.Time
}