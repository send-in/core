package controller

import (
	middleware "core/api/middleware"
	model "core/internal/model"
	service "core/internal/service"
	logger "core/pkg/log"
	"math"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

// GetConnections godoc
//
//	@Summary		Get connections
//	@Description	Get all connections belonging to the authenticated account
//	@Tags			connections
//	@Produce		json
//	@Security		CookieAuth
//	@Param			page	query		int		false	"Page number"	default(1)
//	@Param			limit	query		int		false	"Items per page"	default(20)
//	@Param			sort	query		string	false	"recents|a-z|z-a"	default(recents)
//	@Param			q		query		string	false	"Search by name, company, bio or public id"
//	@Param			ids		query		[]string	false	"Filter by public ids"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/connections [get]
func (c *Controller) GetConnections(context *gin.Context) {
	account := middleware.Account(context)
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(context.DefaultQuery("limit", "20"))
	sort := context.DefaultQuery("sort","recents")
	ids := context.QueryArray("ids")
	q := context.DefaultQuery("q", "")

	if page < 1 { page = 1 }
	if limit < 1 { limit = 20 }
	if limit > 100 { limit = 100 }

	order := "created_at DESC"
	switch sort {
		case "a-z":
			order = "first_name ASC, last_name ASC"
		case "z-a":
			order = "first_name DESC, last_name DESC"
		case "recents":
			order = "created_at DESC"
	}

	query := c.db.
		Model(&model.Connection{}).
		Where("account_id = ?", account.ID)

	if q != "" {
		query = query.Where(
			`public_id ILIKE ?
			OR first_name ILIKE ?
			OR last_name ILIKE ?
			OR company ILIKE ?
			OR bio ILIKE ?`,
			"%"+q+"%",
			"%"+q+"%",
			"%"+q+"%",
			"%"+q+"%",
			"%"+q+"%",
		)
	}

	if len(ids) > 0 {
		query = query.
			Where("public_id IN ?", ids)
	}

	var count int64
	if err := query.
		Count(&count).Error;  err != nil {
		logger.Error(
			"Failed to count connections: %v",
			err,
		)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": "Failed to count connections" },
		)

		return
	}

	var connections []model.Connection

	if err := query.
		Order(order).
		Limit(limit).
		Offset((page - 1) * limit).
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
			"total": int(math.Ceil(
				float64(count) / float64(limit),
			)),
			"count": count,
			"page": page,
			"limit": limit,
			"data": connections,
		},
	)
}

// EnrichConnections godoc
//
//	@Summary		Enrich LinkedIn connections
//	@Description	Queues a background job to scrape and sync LinkedIn connections for the authenticated account
//	@Tags			connections
//	@Produce		json
//	@Security		CookieAuth
//	@Success		202	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Router			/connections    [post]
func (c *Controller) EnrichConnections(
	context *gin.Context,
) {
	account := middleware.Account(context)
	service.EnrichmentJobs <- service.EnrichmentRequest{
		Profile:  account.Profile,
		Agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
		Token: "AQEDATOdT50EsmOqAAABnnSKrvAAAAGemJcy8FYAIOEMH6MoHcItN3Qbuqsl4bHsMs-ikkDtcb4YxiGUSslGsV-KNEwBSohR2wrttoKfHyd0q5WcTr1YDd2zkg-e2EAX02Oq08xDDRW18MMJ7NYIWhuh",
		AccountID: account.ID,
		JSession: "ajax:4580714983183004179",
	}

	context.JSON(
		http.StatusAccepted,
		gin.H{
			"message": "Enrichment jobs queued",
		},
	)
}