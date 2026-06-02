package model

import (
	"time"
	"github.com/google/uuid"
)

type Template struct {
	ID uuid.UUID `gorm:"<-:create;type:uuid;default:gen_random_uuid();primaryKey"`
	AccountID *uuid.UUID `gorm:"not null"`
	
	Value string `gorm:"type:text"`
	Label string

	CreatedAt time.Time
	UpdatedAt time.Time
}