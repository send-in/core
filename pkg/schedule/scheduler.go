package scheduler

import (
	model "core/internal/model"
	logger "core/pkg/log"
	parser "core/pkg/parser"

	"time"

	"gorm.io/gorm"
)

type MessageRequest struct {
	UserAgent  string
	JSession   string
	Token      string
	Message    string
	Receiver   string
	Name       string
	ProfileURN string
	Recipient  string
}


func Start(db *gorm.DB) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	jobs := make(chan MessageRequest, 100)
	go Worker(jobs)

	for {
		<-ticker.C
		var messages []model.Message
		logger.Info("⏳ Tick")

		err := db.
			Preload("Account").
			Preload("Template").
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
			var text string
			if(message.TemplateID != nil){
				if  text, err = parser.ParseLexical(message.Template.Value);
					err != nil {
					logger.Error("Error in parsing template: %s", err)
					continue
				}
			} else if (message.Message != nil) {
				if  text, err = parser.ParseLexical(*message.Message);
				 	err != nil {
					logger.Error("Error in parsing message: %s", err)
					continue
				}
			}

			// TODO: add a fallback
			request := MessageRequest{
				UserAgent: message.Account.UserAgent,
				JSession: message.Account.Session,
				Token: message.Account.Token,
				Receiver: message.Profile,
				Name: message.Name,
				
				ProfileURN: message.ProfileURN,
				Recipient: message.Recipient,
				
				Message: text,
			}

			jobs <- request
		}
	}
}