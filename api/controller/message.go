package controller

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

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

type UpdateMessageRequest struct {
	Timezone     string     `json:"timezone"`
	Message      *string    `json:"message"`
	TemplateID   *uuid.UUID `json:"templateId"`
	ScheduleTime string     `json:"scheduleTime"`
}

// GetMessages godoc
//
//	@Summary		Get messages
//	@Description	Get paginated messages belonging to the authenticated account
//	@Tags			messages
//	@Produce		json
//	@Security		CookieAuth
//	@Param			page	query		int		false	"Page number"	default(1)
//	@Param			limit	query		int		false	"Items per page"	default(20)
//	@Param			sort	query		string	false	"recents|a-z|z-a"	default(recents)
//	@Param			q		query		string	false	"Search by name or company"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/messages [get]
func (c *Controller) GetMessages(context *gin.Context) {
	account := middleware.Account(context)

	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(context.DefaultQuery("limit", "20"))
	sort := context.DefaultQuery("sort","recents")
	q := context.DefaultQuery("q", "")

	if page < 1 { page = 1 }
	if limit < 1 { limit = 20 }
	if limit > 100 { limit = 100 }

	order := "created_at DESC"
	switch sort {
		case "a-z":
			order = "name ASC"
		case "z-a":
			order = "name DESC"
		case "recents":
			order = "created_at DESC"
	}

	query := c.db.
		Model(&model.Message{}).
		Preload("Template").
		Where("account_id = ?", account.ID)

	if q != "" {
		query = query.Where(
			"name ILIKE ? OR company ILIKE ?",
			"%"+q+"%",
			"%"+q+"%",
		)
	}

	var count int64
	if err := query.
		Count(&count).Error;  err != nil {
		logger.Error(
			"Failed to count messages: %v",
			err,
		)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)

		return
	}

	var messages []model.Message

	if err := query.
		Order(order).
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&messages).Error; 
		err != nil {

		logger.Error("Failed to find messages: %v", err)
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"total": int(math.Ceil(
				float64(count) / float64(limit),
			)),
			"count": count,
			"page": page,
			"limit": limit,
			"data": messages,
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
				gin.H{ "error": err.Error() },
			)
			return
		}

		logger.Error("Failed to find message: %v", err)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
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

	var requests []CreateMessageRequest
	if err := context.ShouldBindJSON(&requests); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if len(requests) == 0 {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "No messages provided",
			},
		)
		return
	}

	// if !account.CanSchedule() {
	// 	context.JSON(
	// 		http.StatusForbidden,
	// 		gin.H{
	// 			"error": "Message limit reached",
	// 		},
	// 	)
	// 	return
	// }

	messages := make(
		[]model.Message,
		0,
		len(requests),
	)

	for _, request := range requests {
		scheduleTime, err := time.Parse(
			time.RFC3339,
			request.ScheduleTime,
		)

		if err != nil {
			context.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": err.Error(),
				},
			)
			return
		}

		var connection model.Connection
		if err := c.db.
			Where("public_id = ?", request.Profile).
			First(&connection).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				context.JSON(
					http.StatusBadRequest,
					gin.H{
						"error": "Connection not found. Please Resync.",
					},
				)
				return
			}

			context.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			return
		}

		message := model.Message{
			Name:         request.Name,
			IsSent:       false,
			Picture:      request.Picture,
			Profile:      request.Profile,
			Company:      request.Company,
			Message:      request.Message,
			Timezone:     request.Timezone,
			AccountID:    &account.ID,
			ScheduleTime: scheduleTime,
			ProfileURN:   connection.ProfileURN,
			Recipient:    connection.Recipient,
		}

		if request.TemplateID != nil {
			message.TemplateID = request.TemplateID
		}

		messages = append(
			messages,
			message,
		)
	}

	err := c.db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Create(&messages).Error; err != nil {
				return err
			}

			account.CreditsRemaining -= len(messages)
			account.DailySchedulesUsed += len(messages)
			account.LifetimeMessagesUsed += len(messages)

			return tx.Save(&account).Error
		},
	)

	if err != nil {
		logger.Error(
			"Failed to create messages: %v",
			err,
		)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	context.JSON(
		http.StatusCreated,
		gin.H{
			"message": "Messages created",
			"count": len(messages),
			"data": messages,
		},
	)
}

// UpdateMessage godoc
//
//	@Summary		Update message
//	@Description	Update a message belonging to the authenticated account
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id		path		string					true	"Message ID"
//	@Param			body	body		UpdateMessageRequest	true	"Message payload"
//	@Success		200		{object}	model.Message
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/messages/{id} [put]
func (c *Controller) UpdateMessage(context *gin.Context) {
	account := middleware.Account(context)
	id := context.Param("id")

	var request UpdateMessageRequest

	if  err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{ "error": err.Error() },
		)
		return
	}

	var message model.Message
	if  err := c.db.
		Where(
			"id = ? AND account_id = ?",
			id,
			account.ID,
		).
		First(&message).Error; err != nil {
		context.JSON(
			http.StatusNotFound,
			gin.H{ "error": err.Error() },
		)
		return
	}

	if request.Message != nil {
		message.Message = request.Message
		message.TemplateID = nil
	} else if request.TemplateID != nil {
		message.TemplateID = request.TemplateID
		message.Message = nil
	}

	if request.Timezone != "" {
		message.Timezone = request.Timezone
	}

	if request.ScheduleTime != "" {
		scheduleTime, err := time.Parse(
			time.RFC3339,
			request.ScheduleTime,
		)

		if err != nil {
			context.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": err.Error(),
				},
			)
			return
		}

		scheduleTime = scheduleTime.UTC()
		if !message.ScheduleTime.Equal(scheduleTime) {
			if !account.CanSchedule() {
				context.JSON(
					http.StatusForbidden,
					gin.H{
						"error": "Message limit reached",
					},
				)
				return
			}
		}

		message.ScheduleTime = scheduleTime

		account.ConsumeCredit()
		account.DailySchedulesUsed++
		account.LifetimeMessagesUsed++
		if  err := c.db.Save(&account).Error; 
			err != nil {
			logger.Error(
				"Failed to update account usage: %v",
				err,
			)	
		}
	}

	if 	err := c.db.Save(&message).Error; 
		err != nil {
		logger.Error("Failed to update message: %v", err)
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)
		return
	}

	c.db.
		Preload("Template").
		First(&message, message.ID)

	context.JSON(
		http.StatusOK,
		gin.H{ "data": message },
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
			gin.H{ "error": result.Error.Error() },
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