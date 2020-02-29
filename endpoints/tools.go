package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/AnthonyHewins/poly_regression"
)

const (
	MAX_DEGREE = 5
)

type Data struct{
	MaxDeg uint8     `form:"maxDeg" json:"maxDeg"`
	X      []float64 `form:"x"      json:"x"      binding:"required"`
	Y      []float64 `form:"y"      json:"y"      binding:"required"`
}

func PolynomialRegression(c *gin.Context) {
	var X Data

	if err := c.BindJSON(&X); err == nil {
		if MAX_DEGREE < X.MaxDeg {
			c.JSON(400, gin.H{
				"message": "max degree for polynomial regression is 5",
			})
		} else {
			coef, err := polyfit.PolynomialRegression(X.X, X.Y, int(X.MaxDeg))
			if err == nil {
				c.JSON(200, gin.H{"coef": coef})
			} else {
				c.JSON(500, gin.H{"message": err})
			}
		}
	} else {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	}
}
