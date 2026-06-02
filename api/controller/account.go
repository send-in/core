package controller

import (
	logger "core/pkg/log"
	model "core/internal/model"
	middleware "core/api/middleware"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateAccountRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Profile   string `json:"profile"`
	Picture   string `json:"picture"`
	Timezone  string `json:"timezone"`
	UserAgent string `json:"userAgent"`
}

func (c *Controller) GetAccount(context *gin.Context) {
	account := middleware.Account(context)

	context.JSON(
		http.StatusOK,
		gin.H{
			"data": account,
		},
	)
}

func (c *Controller) UpdateAccount(context *gin.Context) {
	account := middleware.Account(context)
	var request UpdateAccountRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	account.Name = request.Name
	account.Email = request.Email
	account.Profile = request.Profile
	account.Picture = request.Picture
	account.Timezone = request.Timezone
	account.UserAgent = request.UserAgent

	if err := c.db.Save(&account).Error; err != nil {
		logger.Error("Failed to update account: %v", err)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Failed to update account",
			},
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{ 
			"data": account,
		},
	)
}

func (c *Controller) DeleteAccount(context *gin.Context) {
	account := middleware.Account(context)

	if err := c.db.Delete(
		&model.Account{},
		"id = ?",
		account.ID,
	).Error; err != nil {

		logger.Error("Failed to delete account: %v", err)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{ 
				"error": "Failed to delete account",
			},
		)
		return
	}

	context.SetCookie(
		"sendin_auth",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "Account deleted successfully",
		},
	)
}