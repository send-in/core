package model

import (
	"time"
	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID `gorm:"<-:create;type:uuid;default:gen_random_uuid();primaryKey"`

	AccountID *uuid.UUID   `gorm:"not null"`
	Account   *Account     `gorm:"foreignKey:AccountID"`

	TemplateID *uuid.UUID
	Template   *Template   `gorm:"foreignKey:TemplateID"`

	Name string
	Picture string
	Profile string
	Company string
	Timezone string
	Message *string
	ScheduleTime time.Time
	IsSent bool
	CreatedAt time.Time
	UpdatedAt time.Time
}