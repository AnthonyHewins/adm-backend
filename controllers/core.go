package controllers

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"

	"github.com/AnthonyHewins/adm-backend/models"
)

type Routes struct {
	// Base URL
	Base               string `yaml:"base"`

	// Unsecured endpoints
	Polyreg            string `yaml:"polyreg"`
	FeatureEngineering string `yaml:"featureEngineering"`
	Registration       string `yaml:"registration"`
	AcctConfirmation   string `yaml:"acctConfirmation"`
	PasswordReset      string `yaml:"passwordReset"`

	// Secure endpoints
	DcfValuation    string `yaml:"tickerValuation"`
}

const (
	// Register errors
	ERR_EMAIL          = "email"
	ERR_ALREADY_EXISTS = "email-exists"
	ERR_PASSWORD       = "password"
	ERR_LATE           = "late"

	// Authentication/authorization
	ERR_UNAUTHORIZED   = "unauthorized"

	// Tools constants
	MAX_DEGREE   = 5
	MAX_ELEMENTS = 100

	// Tools errors
	ERR_CMD             = "cmd"
	ERR_PARAM           = "param"
	ERR_DEGREE          = "deg"
	ERR_LENGTH          = "length"
	ERR_LENGTH_MISMATCH = "length-mismatch"
	ERR_GENERAL         = "server"
)

func Router(r Routes, privkey, pubkey string) *gin.Engine {
	router := gin.Default()

	jwtMiddleware := genAuthMiddleware(privkey, pubkey)


	// Open endpoints
	unsecured := router.Group(r.Base)
	{
		// New accts
		unsecured.POST(r.Registration,       Register)
		unsecured.GET( r.AcctConfirmation,   AcctConfirmation)

		// Other acct stuff
		unsecured.POST(r.PasswordReset,      PasswordReset)

		// Tools
		unsecured.POST(r.Polyreg,            PolynomialRegression)
		unsecured.POST(r.FeatureEngineering, FeatureEngineering)
	}

	// Auth
	auth := router.Group( fmt.Sprintf("%v%v", r.Base, "/auth") )
	{
		// Unprotected
		auth.POST("/login",         jwtMiddleware.LoginHandler)
		auth.GET( "/refresh_token", jwtMiddleware.RefreshHandler)

		// Begin protected
		auth.Use(jwtMiddleware.MiddlewareFunc())
		auth.GET(r.DcfValuation, DcfValuation)
	}

	return router
}

func connectOrError(c *gin.Context) *gorm.DB {
	db, err := models.Connect()

	if err != nil {
		c.JSON(500, gin.H{
			"error": ERR_GENERAL,
			"message": err.Error(),
		})
		return nil
	}

	return db
}

func forceBind(c *gin.Context, form interface{}) bool {
	if err := c.BindJSON(&form); err != nil {
		c.JSON(400, gin.H{
			"error": ERR_PARAM,
			"message": err.Error(),
		})
		return false
	}

	return true
}
