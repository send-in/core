package model

import (
	"time"
	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID `gorm:"<-:create;type:uuid;default:gen_random_uuid();primaryKey"`
	AccountID *uuid.UUID `gorm:"not null"`
	TemplateID *uuid.UUID

	Name string
	Picture string

	Profile string
	Company string
	Timezone string

	Message *string
	Template Template

	ScheduleTime time.Time

	IsSent bool

	CreatedAt time.Time
	UpdatedAt time.Time
}