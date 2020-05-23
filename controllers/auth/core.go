package auth

import (
	"github.com/AnthonyHewins/adm-backend/controllers/api"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const (
	ErrLate         = "late"
	ErrEmail        = "email"
	ErrEmailTaken   = "email-taken"
	ErrWeakPassword = "password-weak"
	ErrUnauthorized = "unauthorized"
)

func AddRoutes(r *gin.Engine, apiBase string, jwtMiddleware *jwt.GinJWTMiddleware) {
	unsecured := r.Group(apiBase)
	{
		// New accts
		unsecured.POST("/register", api.Endpoint(Register))
		unsecured.GET("/confirm-acct", api.Endpoint(AcctConfirmation))

		unsecured.POST("/reset-password", api.Endpoint(PasswordReset))
		unsecured.POST("/confirm-password-reset", api.Endpoint(ConfirmPwReset))

		unsecured.POST("/login", jwtMiddleware.LoginHandler)
	}

	auth := r.Group(apiBase)
	auth.Use(jwtMiddleware.MiddlewareFunc())
	{
		auth.GET("/refresh_token", jwtMiddleware.RefreshHandler)
	}
}
