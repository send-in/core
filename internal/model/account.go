package model

import (
	"time"
	"github.com/google/uuid"
)

type Plan string
const (
	Free Plan = "free"
	Pro  Plan = "pro"
)

type Account struct {
	ID    uuid.UUID `gorm:"<-:create;type:uuid;default:gen_random_uuid();primaryKey"`
	Email string `gorm:"uniqueIndex"`
	Name  string `gorm:"not null"`

	Timezone string
	Profile  string
	Picture  string

	UserAgent string
	Session   string
	Token 	  string

	Onboarding bool
	
	Payments 	[]Payment
	PlanCredits int `gorm:"default:0"`
	Plan 		Plan `gorm:"default:'free'"`

	CreditsRenewAt   *time.Time
	CreditsRemaining int `gorm:"default:0"`

	LastDailyResetAt 	*time.Time
	DailySchedulesUsed  int `gorm:"default:0"`
	DailySyncsUsed 		int `gorm:"default:0"`

	LifetimeMessagesUsed int `gorm:"default:0"`
	LifetimeSyncsUsed 	 int `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Account) IsFree() bool {
	return a.Plan == Free
}

func (a *Account) ResetDailyLimits() {
	now := time.Now().UTC()

	if a.LastDailyResetAt == nil {
		a.DailySchedulesUsed = 0
		a.DailySyncsUsed = 0
		a.LastDailyResetAt = &now
		return
	}

	if now.YearDay() != a.LastDailyResetAt.YearDay() ||
		now.Year() != a.LastDailyResetAt.Year() {

		a.DailySchedulesUsed = 0
		a.DailySyncsUsed = 0
		a.LastDailyResetAt = &now
	}
}

func (a *Account) RenewCredits() {
	if a.CreditsRenewAt == nil {
		return
	}

	now := time.Now().UTC()
	if !now.After(*a.CreditsRenewAt) {
		return
	}

	for now.After(*a.CreditsRenewAt) {
		next := a.CreditsRenewAt.AddDate(0, 1, 0)
		a.CreditsRenewAt = &next
	}

	a.CreditsRemaining = a.PlanCredits
}

func (a *Account) CanSchedule() bool {
	a.ResetDailyLimits()
	a.RenewCredits()

	if a.IsFree() {
		return a.CreditsRemaining > 0 &&
			a.DailySchedulesUsed < 1 &&
			a.LifetimeMessagesUsed < 5
	}

	return a.CreditsRemaining > 0 &&
		a.DailySchedulesUsed < 10
}

func (a *Account) CanSync() bool {
	a.ResetDailyLimits()
	if a.IsFree() {
		return a.LifetimeSyncsUsed < 1
	}

	return a.DailySyncsUsed < 5
}

func (a *Account) ConsumeCredit() {
	if a.CreditsRemaining > 0 {
		a.CreditsRemaining--
	}
}