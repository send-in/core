package middleware

import (
	"errors"
	"net/http"

	model "core/internal/model"
	logger "core/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)


func Session(context *gin.Context) (string, error) {
	session, err := gothic.Store.Get(
		context.Request,
		"sendin_auth",
	)

	if err != nil { 
		return "", err 
	}

	value, exists := session.Values["account"]
	if !exists {
		return "", errors.New(
			"missing account session",
		)
	}

	accountID, ok := value.(string)
	if !ok {
		return "", errors.New(
			"invalid account session",
		)
	}

	return accountID, nil
}

func Authenticate(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		accountID, err := Session(context)
		if err != nil {
			context.JSON(
				http.StatusUnauthorized,
				gin.H{ "error": err.Error() },
			)
			context.Abort()
			return
		}

		var account model.Account
		if  err := db.
			Where("id = ?", accountID).
			First(&account).Error; 
			err != nil {

			logger.Warning(
				"Failed authentication for account %s",
				accountID,
			)

			context.JSON(
				http.StatusUnauthorized,
				gin.H{ "error": err.Error() },
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