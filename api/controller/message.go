package controller

import (
	middleware "core/api/middleware"
	model "core/internal/model"
	logger "core/pkg/log"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateMessageRequest struct {
	Name          string      `json:"name" binding:"required"`
	Picture       string      `json:"picture"`
	Profile       string      `json:"profile" binding:"required"`
	Company       string      `json:"company"`
	Timezone      string      `json:"timezone"`
	Message       *string     `json:"message"`
	TemplateID    *uuid.UUID  `json:"templateId"`
	ScheduleTime  string      `json:"scheduleTime"`
}

func (c *Controller) GetMessages(context *gin.Context) {
	account := middleware.Account(context)

	var messages []model.Message

	if err := c.db.
		Preload("Template").
		Where("account_id = ?", account.ID).
		Find(&messages).Error; err != nil {

		logger.Error("Failed to find messages: %v", err)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to find messages"},
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"count": len(messages),
			"data":  messages,
		},
	)
}

func (c *Controller) GetMessage(context *gin.Context) {
	account := middleware.Account(context)
	id := context.Param("id")

	var message model.Message

	if err := c.db.
		Preload("Template").
		Where(
			"id = ? AND account_id = ?",
			id,
			account.ID,
		).
		First(&message).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(
				http.StatusNotFound,
				gin.H{"error": "Message not found"},
			)
			return
		}

		logger.Error("Failed to find message: %v", err)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to find message"},
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{"data": message},
	)
}

func (c *Controller) CreateMessage(context *gin.Context) {
	account := middleware.Account(context)
	var request CreateMessageRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	message := model.Message{
		Name:      request.Name,
		IsSent:    false,
		Picture:   request.Picture,
		Profile:   request.Profile,
		Company:   request.Company,
		Message:   request.Message,
		Timezone:  request.Timezone,
		AccountID: &account.ID,
	}

	if err := c.db.Create(&message).Error; err != nil {
		logger.Error("Failed to create message: %v", err)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Failed to create message",
			},
		)
		return
	}

	context.JSON(
		http.StatusCreated,
		gin.H{
			"message": "Message created",
			"data":    message,
		},
	)
}

func (c *Controller) DeleteMessage(context *gin.Context) {
	account := middleware.Account(context)
	id := context.Param("id")

	result := c.db.
		Where(
			"id = ? AND account_id = ?",
			id,
			account.ID,
		).
		Delete(&model.Message{})

	if result.Error != nil {
		logger.Error("Failed to delete message: %v", result.Error)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Failed to delete message",
			},
		)
		return
	}

	if result.RowsAffected == 0 {
		context.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "Message not found",
			},
		)

		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "Message deleted successfully",
		},
	)
}