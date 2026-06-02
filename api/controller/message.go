package controller

import (
	"errors"
	"net/http"

	middleware "core/api/middleware"
	model "core/internal/model"
	logger "core/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateMessageRequest struct {
	Name         string     `json:"name" binding:"required"`
	Picture      string     `json:"picture"`
	Profile      string     `json:"profile" binding:"required"`
	Company      string     `json:"company"`
	Timezone     string     `json:"timezone"`
	Message      *string    `json:"message"`
	TemplateID   *uuid.UUID `json:"templateId"`
	ScheduleTime string     `json:"scheduleTime"`
}

// GetMessages godoc
//
//	@Summary		Get messages
//	@Description	Get all messages belonging to the authenticated account
//	@Tags			messages
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{array}		model.Message
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/messages [get]
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
			gin.H{
				"error": "Failed to find messages",
			},
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

// GetMessage godoc
//
//	@Summary		Get message
//	@Description	Get a single message belonging to the authenticated account
//	@Tags			messages
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"Message ID"
//	@Success		200	{object}	model.Message
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/messages/{id} [get]
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
				gin.H{
					"error": "Message not found",
				},
			)
			return
		}

		logger.Error("Failed to find message: %v", err)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Failed to find message",
			},
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"data": message,
		},
	)
}

// CreateMessage godoc
//
//	@Summary		Create message
//	@Description	Create a new message for the authenticated account
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		CreateMessageRequest	true	"Message payload"
//	@Success		201		{object}	model.Message
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/messages [post]
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

	if request.TemplateID != nil {
		message.TemplateID = request.TemplateID
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

// DeleteMessage godoc
//
//	@Summary		Delete message
//	@Description	Delete a message belonging to the authenticated account
//	@Tags			messages
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"Message ID"
//	@Success		200	{object}	map[string]string
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/messages/{id} [delete]
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