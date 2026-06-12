package controller

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	middleware "core/api/middleware"
	model "core/internal/model"
	logger "core/pkg/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTemplateRequest struct {
	Label string `json:"label" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type UpdateTemplateRequest struct {
	Label string `json:"label" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// GetTemplates godoc
//
//	@Summary		Get templates
//	@Description	Get paginated templates belonging to the authenticated account
//	@Tags			templates
//	@Produce		json
//	@Security		CookieAuth
//	@Param			page	query		int		false	"Page number"	default(1)
//	@Param			limit	query		int		false	"Items per page"	default(20)
//	@Param			sort	query		string	false	"recents|a-z|z-a"	default(recents)
//	@Success		200		{array}		model.Template
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/templates [get]
func (c *Controller) GetTemplates(context *gin.Context) {
	
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
			order = "label ASC"
		case "z-a":
			order = "label DESC"
		case "recents":
			order = "created_at DESC"
	}

	query := c.db.
		Model(&model.Template{}).
		Where("account_id = ?", account.ID)

	if q != "" {
		query = query.Where(
			"label ILIKE ? OR value ILIKE ?",
			"%"+q+"%",
			"%"+q+"%",
		)
	}

	var count int64
	if err := query.
		Count(&count).Error; 
		err != nil {
		logger.Error(
			"Failed to count templates: %v",
			err,
		)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)

		return
	}

	var templates []model.Template

	if  err := query.
		Order(order).
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&templates).Error;
		err != nil {

		logger.Error(
			"Failed to find templates: %v",
			err,
		)

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
			"data": templates,
		},
	)
}

// GetTemplate godoc
//
//	@Summary		Get template
//	@Description	Get a single template belonging to the authenticated account
//	@Tags			templates
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"Template ID"
//	@Success		200	{object}	model.Template
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/templates/{id} [get]
func (c *Controller) GetTemplate(context *gin.Context) {
	account := middleware.Account(context)
	id := context.Param("id")

	var template model.Template

	if err := c.db.
		Where(
			"id = ? AND account_id = ?",
			id,
			account.ID,
		).
		First(&template).Error; err != nil {

		if  errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(
				http.StatusNotFound,
				gin.H{ "error": err.Error() },
			)
			return
		}

		logger.Error("Failed to find template: %v", err)
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"data": template,
		},
	)
}

// CreateTemplate godoc
//
//	@Summary		Create template
//	@Description	Create a new template for the authenticated account
//	@Tags			templates
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		CreateTemplateRequest	true	"Template payload"
//	@Success		201		{object}	model.Template
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/templates [post]
func (c *Controller) CreateTemplate(context *gin.Context) {
	account := middleware.Account(context)
	if account.IsFree() {
		var count int64

		if err := c.db.
			Model(&model.Template{}).
			Where(
				"account_id = ?",
				account.ID,
			).
			Count(&count).Error; err != nil {

			logger.Error(
				"Failed to count templates: %v",
				err,
			)

			context.JSON(
				http.StatusInternalServerError,
				gin.H{ "error": err.Error() },
			)

			return
		}

		if count >= 1 {
			context.JSON(
				http.StatusForbidden,
				gin.H{
					"error": "Free plan allows only 1 template",
				},
			)

			return
		}
	}

	var request CreateTemplateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	template := model.Template{
		AccountID: &account.ID,
		Label:     request.Label,
		Value:     request.Value,
	}

	if err := c.db.Create(&template).Error; err != nil {
		logger.Error("Failed to create template: %v", err)
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)
		return
	}

	context.JSON(
		http.StatusCreated,
		gin.H{
			"message": "Template created",
			"data":    template,
		},
	)
}

// UpdateTemplate godoc
//
//	@Summary		Update template
//	@Description	Update a template belonging to the authenticated account
//	@Tags			templates
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id		path		string					true	"Template ID"
//	@Param			body	body		UpdateTemplateRequest	true	"Template payload"
//	@Success		200		{object}	model.Template
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/templates/{id} [put]
func (c *Controller) UpdateTemplate(context *gin.Context) {
	account := middleware.Account(context)
	id := context.Param("id")

	var request UpdateTemplateRequest
	if  err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{ "error": err.Error() },
		)
		return
	}

	var template model.Template

	if err := c.db.
		Where(
			"id = ? AND account_id = ?",
			id,
			account.ID,
		).
		First(&template).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(
				http.StatusNotFound,
				gin.H{ "error": err.Error() },
			)
			return
		}

		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)
		return
	}

	template.Label = request.Label
	template.Value = request.Value

	if  err := c.db.Save(&template).Error; 
	 	err != nil {
		logger.Error("Failed to update template: %v", err)
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{ "data": template },
	)
}

// DeleteTemplate godoc
//
//	@Summary		Delete template
//	@Description	Delete a template belonging to the authenticated account
//	@Tags			templates
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"Template ID"
//	@Success		200	{object}	map[string]string
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/templates/{id} [delete]
func (c *Controller) DeleteTemplate(context *gin.Context) {
	account := middleware.Account(context)
	id := context.Param("id")

	result := c.db.
		Where(
			"id = ? AND account_id = ?",
			id,
			account.ID,
		).
		Delete(&model.Template{})

	if  result.Error != nil {
		logger.Error("Failed to delete template: %v", result.Error)
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": result.Error.Error() },
		)
		return
	}

	if  result.RowsAffected == 0 {
		context.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "Template not found",
			},
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "Template deleted successfully",
		},
	)
}