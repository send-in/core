package controller

import (
	"core/api/middleware"
	model "core/internal/model"
	logger "core/pkg/log"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CreateConnectionsRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

func (c *Controller) GetConnections(context *gin.Context) {
	account := middleware.Account(context)
	var connections []model.Connection

	if err := c.db.
		Where("account_id = ?", account.ID).
		Find(&connections).Error; err != nil {

		logger.Error("Failed to find connections: %v", err)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Failed to find connections",
			},
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"count": len(connections),
			"data":  connections,
		},
	)
}

func (c *Controller) CreateConnections(context *gin.Context) {
	account := middleware.Account(context)

	var request CreateConnectionsRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	var connections []model.Connection

	for _, publicID := range request.IDs {
		publicID = strings.TrimSpace(publicID)

		if publicID == "" {
			continue
		}

		var connection model.Connection

		err := c.db.
			Where(
				"account_id = ? AND public_id = ?",
				account.ID,
				publicID,
			).
			First(&connection).
			Error

		if err == nil {
			connections = append(connections, connection)
			continue
		}

		connection = model.Connection{
			AccountID: &account.ID,
			PublicID:  publicID,
		}

		if err := c.db.Create(&connection).Error; err != nil {
			logger.Error(
				"Failed to create connection %s: %v",
				publicID,
				err,
			)
			continue
		}

		connections = append(connections, connection)
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "Connections created successfully",
			"data":    connections,
		},
	)
}