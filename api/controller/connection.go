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

// GetConnections godoc
//
//	@Summary		Get connections
//	@Description	Get all connections belonging to the authenticated account
//	@Tags			connections
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{array}		model.Connection
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/connections [get]
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

// CreateConnections godoc
//
//	@Summary		Create connections
//	@Description	Create or retrieve connections from a list of LinkedIn public IDs
//	@Tags			connections
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		CreateConnectionsRequest	true	"Connection IDs"
//	@Success		200		{array}		model.Connection
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/connections [post]
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