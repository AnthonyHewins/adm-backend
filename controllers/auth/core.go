package auth

import (
	"fmt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const (
	ErrLate = "late"
)

func AddRoutes(r *gin.Engine, apiBase string, jwtMiddleware *jwt.GinJWTMiddleware) {
	unsecured := r.Group(apiBase)
	{
		// New accts
		unsecured.POST("/register", Register)
		unsecured.GET( "/confirm-acct", AcctConfirmation)

		// Active account actions
		unsecured.POST("/reset-password", PasswordReset)
		unsecured.POST("/confirm-reset", ConfirmPwReset)
	}

	auth := r.Group( fmt.Sprintf("%v%v", apiBase, "/auth") )
	auth.Use(jwtMiddleware.MiddlewareFunc())
	{
		// Unprotected
		auth.POST("/login",         jwtMiddleware.LoginHandler)
		auth.GET( "/refresh_token", jwtMiddleware.RefreshHandler)
	}
}
