package scheduler

import (
	"bytes"
	"encoding/json"
	"net/http"

	logger "core/pkg/log"
)

func Worker(jobs <-chan MessageRequest) {
	for message := range jobs {
		bypass := false

		logger.Info(
			"📤 Dispatching scheduled message to %s",
			message.Receiver,
		)

		body, err := json.Marshal(message)
		if err != nil {
			logger.Error(
				"Failed to marshal message to %s: %v",
				message.Receiver,
				err,
			)
			return
		}

		if(bypass){
			logger.Success(
				"✅ Scheduled message to %s dispatched",
				message.Receiver,
			)
			return
		}

		response, err := http.Post(
			"http://localhost:8001/api/v1/jobs",
			"application/json",
			bytes.NewBuffer(body),
		)

		if err != nil {
			logger.Error(
				"❌ Failed to dispatch message to %s: %v",
				message.Receiver,
				err,
			)
			return
		}

		defer response.Body.Close()
		if response.StatusCode >= 300 {
			logger.Error(
				"🚫 Consumer rejected message to %s with status %d",
				message.Receiver,
				response.StatusCode,
			)
			return
		}

		logger.Success(
			"✅ Scheduled message to %s dispatched",
			message.Receiver,
		)
	}
}