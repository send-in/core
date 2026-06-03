package scheduler

import (
	"bytes"
	"encoding/json"
	"net/http"

	model "core/internal/model"
	logger "core/pkg/log"
)

func Worker(message model.Message) {
	bypass := true

	logger.Info(
		"📤 Dispatching scheduled message %s",
		message.ID,
	)

	body, err := json.Marshal(message)
	if err != nil {
		logger.Error(
			"Failed to marshal message %s: %v",
			message.ID,
			err,
		)
		return
	}

	if(bypass){
		logger.Success(
			"✅ Scheduled message %s dispatched",
			message.ID,
		)
		return
	}

	response, err := http.Post(
		"http://localhost:8001/jobs",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		logger.Error(
			"❌ Failed to dispatch message %s: %v",
			message.ID,
			err,
		)
		return
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		logger.Error(
			"🚫 Consumer rejected message %s with status %d",
			message.ID,
			response.StatusCode,
		)
		return
	}

	logger.Success(
		"✅ Scheduled message %s dispatched",
		message.ID,
	)
}