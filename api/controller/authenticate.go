package controller

import (
	"errors"
	"net/http"

	model "core/internal/model"
	logger "core/pkg/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SignupRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Profile   string `json:"profile"`
	Picture   string `json:"picture"`
	Timezone  string `json:"timezone"`
	Token     string `json:"token"`
	UserAgent string `json:"userAgent"`
}

func (c *Controller) Login(context *gin.Context) {
	var request LoginRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	var account model.Account

	if err := c.db.
		Where("email = ?", request.Email).
		First(&account).Error; err != nil {

		context.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Invalid credentials"},
		)
		return
	}

	context.SetCookie(
		"sendin_auth",
		account.ID.String(),
		3600*24*30,
		"/",
		"",
		false,
		true,
	)

	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "Login successful",
			"data":    account,
		},
	)
}

func (c *Controller) Signup(context *gin.Context) {
	var request SignupRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	var existing model.Account

	err := c.db.
		Where("email = ?", request.Email).
		First(&existing).
		Error

	if err == nil {
		context.JSON(
			http.StatusConflict,
			gin.H{"error": "Account already exists"},
		)
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to check account existence: %v", err)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create account"},
		)
		return
	}

	account := model.Account{
		Name:      request.Name,
		Email:     request.Email,
		Profile:   request.Profile,
		Picture:   request.Picture,
		Timezone:  request.Timezone,
		Token:     request.Token,
		UserAgent: request.UserAgent,
	}

	if err := c.db.Create(&account).Error; err != nil {
		logger.Error("Failed to create account: %v", err)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create account"},
		)
		return
	}

	context.SetCookie(
		"sendin_auth",
		account.ID.String(),
		3600*24*30,
		"/",
		"",
		false,
		true,
	)

	context.JSON(
		http.StatusCreated,
		gin.H{
			"message": "Account created",
			"data":    account,
		},
	)
}