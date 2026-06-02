package middleware

import (
	"net/http"

	model "core/internal/model"
	logger "core/pkg/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Authenticate(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		accountID, err := context.Cookie("sendin_auth")

		if err != nil {
			context.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Missing authentication cookie",
				},
			)
			context.Abort()
			return
		}

		var account model.Account

		if err := db.
			Where("id = ?", accountID).
			First(&account).Error; err != nil {

			logger.Warning(
				"Failed authentication for account %s",
				accountID,
			)

			context.JSON(
				http.StatusUnauthorized,
				gin.H{ "error": "Invalid session" },
			)

			context.Abort()
			return
		}

		context.Set("account", account)
		context.Next()
	}
}

func Account(context *gin.Context) model.Account {
	return context.MustGet("account").(model.Account)
}