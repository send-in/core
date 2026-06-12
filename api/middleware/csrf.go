package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func CSRF() gin.HandlerFunc {
	return func(context *gin.Context) {

		if context.Request.Method == http.MethodGet {
			context.Next()
			return
		}

		session, err := gothic.Store.Get(
			context.Request,
			"sendin_auth",
		)

		if err != nil {
			context.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "Invalid session",
				},
			)
			return
		}

		expected, ok :=
			session.Values["csrf"].(string)

		if !ok {
			context.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "Missing CSRF token",
				},
			)
			return
		}

		actual :=
			context.GetHeader(
				"X-CSRF-Token",
			)

		if actual != expected {
			context.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "Invalid CSRF token",
				},
			)
			return
		}

		context.Next()
	}
}