package model

import (
	"time"
	"github.com/google/uuid"
)

type Connection struct {
	ID 		  uuid.UUID `gorm:"<-:create;type:uuid;default:gen_random_uuid();primaryKey"`
	PublicID  string `gorm:"uniqueIndex"`

	FirstName string
	LastName  string

	Bio		string
	Picture string

	Company   string
	Country   string
	Timezone *string

	ProfileURN string
	Recipient  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountConnection struct {
    AccountID    uuid.UUID `gorm:"primaryKey"`
	Account      *Account  `gorm:"foreignKey:AccountID"`

    ConnectionID uuid.UUID   `gorm:"primaryKey"`
	Connection   *Connection `gorm:"foreignKey:ConnectionID"`

    CreatedAt time.Time
}