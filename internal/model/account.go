package model

import (
	"time"
	"github.com/google/uuid"
)

type Account struct {
	ID uuid.UUID `gorm:"<-:create;type:uuid;default:gen_random_uuid();primaryKey"`
	Name string `gorm:"not null"`
	Email string `gorm:"uniqueIndex"`

	Profile string
	Picture string
	Timezone string

	Token string
	UserAgent string 

	CreatedAt time.Time
	UpdatedAt time.Time
}