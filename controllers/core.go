package controllers

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"

	"github.com/AnthonyHewins/adm-backend/models"
)

type Routes struct {
	Polyreg            string `yaml:"polyreg"`
	FeatureEngineering string `yaml:"featureEngineering"`
	Registration       string `yaml:"registration"`
	AcctConfirmation   string `yaml:"acctConfirmation"`
}

const (
	// Register errors
	ERR_EMAIL    = "email"
	ERR_PASSWORD = "password"
	ERR_LATE     = "late"

	// Tools constants
	MAX_DEGREE   = 5
	MAX_ELEMENTS = 100

	// Tools errors
	ERR_CMD              = "cmd"
	ERR_PARAM            = "param"
	ERR_DEGREE           = "deg"
	ERR_LENGTH           = "length"
	ERR_LENGTH_MISMATCH  = "length-mismatch"
	ERR_GENERAL          = "server"
)

func Router(r *Routes) *gin.Engine {
	router := gin.Default()

	router.POST(r.Registration,       Register)
	router.GET( r.AcctConfirmation,   AcctConfirmation)

	router.POST(r.Polyreg,            PolynomialRegression)
	router.POST(r.FeatureEngineering, FeatureEngineering)

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
