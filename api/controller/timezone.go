package controller

import (
	openai "core/pkg/openai"
	"net/http"
	"github.com/gin-gonic/gin"
)

type InferTimezoneRequest struct {
    Location string `json:"location" binding:"required"`
}

type InferTimezoneResponse struct {
    Country  string `json:"country"`
    Timezone string `json:"timezone"`
}

// InferTimezone godoc
//
//	@Summary		Infer timezone from a LinkedIn profile
//	@Description	Infers country and IANA timezone from a parsed LinkedIn profile location
//	@Tags			connections
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		InferTimezoneRequest	true	"Profile"
//	@Success		200		{object}	InferTimezoneResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/timezone [post]
func (c *Controller) InferTimezone(context *gin.Context) {
	var request InferTimezoneRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	result, err := openai.Client.InferTimezone(
		request.Location,
	)

	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"country": result.Country,
			"timezone": result.Timezone,
		},
	)
}