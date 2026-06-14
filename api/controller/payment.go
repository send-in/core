package controller

import (
	logger "core/pkg/log"
	model "core/internal/model"
	payment "core/internal/payment"
	middleware "core/api/middleware"

	"time"
	"errors"
	"strconv"
	"net/http"
	"encoding/json"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePaymentRequest struct {
	Credits int `json:"credits" binding:"required,min=25,max=200"`
}

type RazorpayWebhook struct {
	Event string `json:"event"`

	Payload struct {
		Payment struct {
			Entity struct {
				ID      string `json:"id"`
				OrderID string `json:"order_id"`
			} `json:"entity"`
		} `json:"payment"`
	} `json:"payload"`
}

// GetPayment godoc
//
//	@Summary		Get payment
//	@Description	Get payment status by Razorpay order id
//	@Tags			payment
//	@Produce		json
//	@Security		CookieAuth
//	@Router			/payments/{orderId} [get]
func (c *Controller) GetPayment(context *gin.Context) {
	account := middleware.Account(context)
	orderID := context.Param("orderId")

	var payment model.Payment
	if err := c.db.
		Where(
			"account_id = ? AND order_id = ?",
			account.ID,
			orderID,
		).First(&payment).Error; 
		err != nil {
		context.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "Payment not found",
			},
		)

		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"data": gin.H{
				"status": payment.Status,
				"credits": payment.PlanCredits,
			},
		},
	)
}

// CreatePayment godoc
//
//	@Summary		Create payment order
//	@Description	Creates a Razorpay order for purchasing SendIn credits
//	@Tags			payment
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		CreatePaymentRequest	true	"Payment request"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		409		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/payments [post]
func (c *Controller) CreatePayment(context *gin.Context) {
	account := middleware.Account(context)

	var request CreatePaymentRequest
	if err := context.
		ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	if 	request.Credits < payment.MinCredits ||
		request.Credits > payment.MaxCredits {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Credits must be between 25 and 200",
			},
		)
		return
	}

	if request.Credits % payment.CreditStep != 0 {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Credits must be purchased in increments of 25",
			},
		)

		return
	}

	amount := int64(
		request.Credits *
		payment.PricePerCredit *
		100,
	)

	order, err := payment.CreateOrder(
		payment.OrderRequest{
			Amount:   amount,
			Currency: "USD",
			Receipt:  uuid.NewString(),
			Notes: map[string]string{
				"accountId": account.ID.String(),
				"credits":   strconv.Itoa(request.Credits),
				"plan":      "pro",
			},
		},
	)

	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)

		return
	}

	var record model.Payment
	err = c.db.
		Where(
			"account_id = ? AND status = ?",
			account.ID,
			model.PaymentPending,
		).
		Order("created_at DESC").
		First(&record).
		Error

	if err == nil {
		record.Plan = "pro"
		record.PlanCredits = request.Credits
		record.Amount = amount
		record.Currency = "USD"
		record.Provider = "razorpay"
		record.OrderID = order.ID
		record.ExternalID = ""
		record.Payload = nil
		record.CompletedAt = nil

		if  err := c.db.Save(&record).Error; 
			err != nil {

			context.JSON(
				http.StatusInternalServerError,
				gin.H{ "error": err.Error() },
			)

			return
		}
	} else if errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {

		record = model.Payment{
			AccountID: account.ID,
			Status: model.PaymentPending,
			Plan: "pro",
			PlanCredits: request.Credits,
			Amount: amount,
			Currency: "USD",
			Provider: "razorpay",
			OrderID: order.ID,
		}

		if  err := c.db.Create(&record).Error; 
			err != nil {
			context.JSON(
				http.StatusInternalServerError,
				gin.H{ "error": err.Error() },
			)

			return
		}
	} else {
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)

		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"data": gin.H{
				"orderId": order.ID,
				"amount": order.Amount,
				"currency": order.Currency,
			},
		},
	)
}

// RazorpayWebhook godoc
//
//	@Summary		Razorpay webhook
//	@Description	Processes successful Razorpay payments
//	@Tags			payment
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Failure		401
//	@Router			/payments/webhook [post]
func (c *Controller) RazorpayWebhook(context *gin.Context) {
	body, err := context.GetRawData()
	if err != nil {
		context.Status(
			http.StatusBadRequest,
		)
		return
	}

	signature := context.GetHeader("X-Razorpay-Signature")
	if !payment.VerifyWebhook(
		body,
		signature,
	) {
		context.Status(
			http.StatusUnauthorized,
		)
		return
	}

	var webhook RazorpayWebhook
	if  err := json.Unmarshal(body, &webhook); 
		err != nil {
		context.Status(
			http.StatusBadRequest,
		)
		return
	}

	if webhook.Event != "payment.captured" {
		context.Status(http.StatusOK)
		return
	}

	var record model.Payment
	if err := c.db.
		Where(
			"order_id = ?",
			webhook.Payload.
				Payment.
				Entity.
				OrderID,
		).First(&record).Error; 
		err != nil {

		logger.Warning(
			"Payment not found for order %s",
			webhook.Payload.Payment.Entity.OrderID,
		)

		context.Status(http.StatusNotFound)
		return
	}

	if record.ExternalID != "" {
		context.Status(http.StatusOK)
		return
	}

	if record.Status ==
		model.PaymentSucceeded {
		context.Status(http.StatusOK)
		return
	}

	var account model.Account
	if  err := c.db.
		First(&account, record.AccountID).
		Error; 
		err != nil {

		context.Status(http.StatusNotFound)
		return
	}

	now := time.Now()

	err = c.db.Transaction(
		func(tx *gorm.DB) error {
			record.Status = model.PaymentSucceeded
			record.ExternalID = webhook.Payload.Payment.Entity.ID
			record.Payload = body
			record.CompletedAt = &now

			account.Plan = model.Pro
			account.PlanCredits = record.PlanCredits

			if account.CreditsRemaining < record.PlanCredits {
				account.CreditsRemaining = record.PlanCredits
			}

			if account.CreditsRenewAt == nil {
				next := now.AddDate(0, 1, 0)
				account.CreditsRenewAt = &next
			}

			if 	err := tx.Save(&record).Error; 
				err != nil {
				return err
			}

			if  err := tx.Save(&account).Error; 
				err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		logger.Error(
			"Failed processing payment %s: %v",
			record.ID,
			err,
		)

		context.Status(http.StatusInternalServerError)
		return
	}

	logger.Info(
		"Payment %s succeeded for account %s",
		record.ID,
		account.ID,
	)

	context.Status(http.StatusOK)
}

// CancelSubscription godoc
//
//	@Summary		Cancel subscription
//	@Description	Cancels the current Pro subscription while preserving any remaining credits
//	@Tags			payment
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/payments/subscription [delete]
func (c *Controller) CancelSubscription(
	context *gin.Context,
) {
	account := middleware.Account(context)

	if account.Plan == model.Free {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "No active subscription",
			},
		)
		return
	}

	account.Plan = model.Free
	account.PlanCredits = 0
	account.CreditsRenewAt = nil

	if err := c.db.
		Save(&account).
		Error; err != nil {

		logger.Error(
			"Failed to cancel subscription for account %s: %v",
			account.ID,
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
		http.StatusOK,
		gin.H{
			"message": "Subscription cancelled",
			"data": gin.H{
				"plan": account.Plan,
				"creditsRemaining": account.CreditsRemaining,
			},
		},
	)
}