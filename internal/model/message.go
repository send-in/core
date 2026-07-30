package model

import (
	"time"
	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID `gorm:"<-:create;type:uuid;default:gen_random_uuid();primaryKey"`

	AccountID 	*uuid.UUID   `gorm:"not null"`
	Account   	*Account     `gorm:"foreignKey:AccountID"`

	TemplateID  *uuid.UUID
	Template    *Template    `gorm:"foreignKey:TemplateID"`
	Message 	*string

	Name 		string
	Picture 	string
	Profile 	string
	Company 	string
	Timezone 	string
	Recipient 	string
	ProfileURN  string
	
	IsSent 		 bool
	CreatedAt 	 time.Time
	UpdatedAt 	 time.Time
	ScheduleTime time.Time
}