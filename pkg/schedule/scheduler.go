package scheduler

import (
	model "core/internal/model"
	logger "core/pkg/log"

	"time"
	"gorm.io/gorm"
)

func Start(db *gorm.DB) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	jobs := make(chan model.Message, 100)
	go Worker(jobs)

	for {
		<-ticker.C
		var messages []model.Message
		logger.Info("⏳ Tick")

		err := db.
			Where("is_sent = ?", false).
			Where("schedule_time <= ?", time.Now()).
			Find(&messages).
			Error

		if err != nil {
			logger.Error(
				"Failed to fetch scheduled messages: %v",
				err,
			)
			continue
		}

		for _, message := range messages {
			jobs <- message
		}
	}
}