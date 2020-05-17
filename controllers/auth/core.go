package auth

import (
	"fmt"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const (
	ErrLate = "late"
	ErrEmail = "email"
	ErrEmailTaken = "email-taken"
	ErrWeakPassword = "password-weak"
	ErrUnauthorized = "unauthorized"
)

func AddRoutes(r *gin.Engine, apiBase string, jwtMiddleware *jwt.GinJWTMiddleware) {
	unsecured := r.Group(apiBase)
	{
		// New accts
		unsecured.POST("/register", api.Endpoint(Register))
		unsecured.GET( "/confirm-acct", api.Endpoint(AcctConfirmation))

		// Active account actions
		//unsecured.POST("/reset-password", api.Endpoint(PasswordReset))
		//unsecured.POST("/confirm-reset", api.Endpoint(ConfirmPwReset))
	}

	auth := r.Group( fmt.Sprintf("%v%v", apiBase, "/auth") )
	auth.Use(jwtMiddleware.MiddlewareFunc())
	{
		// Unprotected
		auth.POST("/login",         jwtMiddleware.LoginHandler)
		auth.GET( "/refresh_token", jwtMiddleware.RefreshHandler)
	}
}
