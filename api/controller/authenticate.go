package controller

import (
	"errors"
	"net/http"
	"strconv"

	model "core/internal/model"
	logger "core/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

// LinkedInLogin godoc
//
//	@Summary		Login with LinkedIn
//	@Description	Redirects the user to LinkedIn OAuth authentication
//	@Tags			auth
//	@Produce		json
//	@Success		307
//	@Router			/auth/linkedin [get]
func (c *Controller) LinkedInLogin(context *gin.Context) {
	query := context.Request.URL.Query()
	query.Set("provider", "linkedin-oidc")
	context.Request.URL.RawQuery = query.Encode()

	gothic.BeginAuthHandler(
		context.Writer,
		context.Request,
	)
}

// LinkedInCallback godoc
//
//	@Summary		LinkedIn callback
//	@Description	Handles LinkedIn OAuth callback, creates account if needed, authenticates user and sets session cookie
//	@Tags			auth
//	@Produce		json
//	@Success		307
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/auth/linkedin/callback [get]
func (c *Controller) LinkedInCallback(context *gin.Context) {
	query := context.Request.URL.Query()
	query.Set("provider", "linkedin-oidc")
	context.Request.URL.RawQuery = query.Encode()
	
	user, err := gothic.CompleteUserAuth(
		context.Writer,
		context.Request,
	)

	if err != nil {
		logger.Error("LinkedIn authentication failed: %v", err)
		context.JSON(
			http.StatusUnauthorized,
			gin.H{ "error": err.Error() },
		)
		return
	}

	var account model.Account
	err = c.db.
		Where(
			"email = ?",
			user.Email,
		).
		First(&account).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		account = model.Account{
			Name: user.Name,
			Email: user.Email,
			Profile: user.NickName,
			Picture: user.AvatarURL,
			Plan: model.Free,
			PlanCredits: 5,
			CreditsRemaining: 5,
			DailySchedulesUsed: 0,
			DailySyncsUsed: 0,
			LifetimeSyncsUsed: 0,
			LifetimeMessagesUsed: 0,
			Onboarding: false,
		}

		if 	err := c.db.
			Create(&account).
			Error; 
			err != nil {
			logger.Error("Failed to create account: %v", err)
			context.JSON(
				http.StatusInternalServerError,
				gin.H{ "error": err.Error() },
			)
			return
		}
	} else if err != nil {
		logger.Error("Failed to lookup account: %v", err )
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)

		return
	}

	session, err := gothic.Store.Get(
		context.Request,
		"sendin_auth",
	)

	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)
		return
	}

	session.Values["account"] = account.ID.String()
	if  err := session.Save(context.Request, context.Writer);  
		err != nil {
		context.JSON(
			http.StatusInternalServerError,
			gin.H{ "error": err.Error() },
		)
		return
	}

	http.SetCookie(
		context.Writer,
		&http.Cookie{
			Name:     "sendin_onboarded",
			Value:    strconv.FormatBool(account.Onboarding),
			Path:     "/",
			HttpOnly: false,
		},
	)

	http.SetCookie(
		context.Writer, 
		&http.Cookie{
			Name:  "sendin_access",
			Value: "true",
			Path:  "/",
			MaxAge:   86400 * 30,
			HttpOnly: false,
		},
	)

	context.Redirect(
		http.StatusTemporaryRedirect,
		"http://localhost:3000",
	)
}


// Logout godoc
//
//	@Summary		Logout
//	@Description	Clears the current authentication session
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/auth/logout [post]
func (c *Controller) Logout(context *gin.Context) {
	session, _ := gothic.Store.Get(
		context.Request,
		"sendin_auth",
	)

	session.Options.MaxAge = -1
	_ = session.Save(context.Request, context.Writer)

	http.SetCookie(
		context.Writer,
		&http.Cookie{
			Name:     "sendin_access",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: false,
		},
	)

	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "Logged out",
		},
	)
}
