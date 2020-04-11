package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/AnthonyHewins/polyfit"
)

type PolyRegData struct{
	MaxDeg *int       `form:"maxDeg" json:"maxDeg" binding:"required"`
	X      *[]float64 `form:"x"      json:"x"      binding:"required"`
	Y      *[]float64 `form:"y"      json:"y"      binding:"required"`
}

func PolynomialRegression(c *gin.Context) {
	var X PolyRegData

	if err := c.BindJSON(&X); err != nil {
		c.JSON(400, gin.H{
			"error": ERR_PARAM,
			"message": err.Error(),
		})
		return
	}

	if MAX_DEGREE < *X.MaxDeg || 0 >= *X.MaxDeg {
		c.JSON(422, gin.H{
			"error": ERR_DEGREE,
			"message": fmt.Sprintf("maxDeg must satisfy 0 <= maxDeg <= %v", MAX_DEGREE),
		})
		return
	}

	n := len(*X.X)
	m := len(*X.Y)

	if (n != m || n > MAX_ELEMENTS || n <= *X.MaxDeg) {
		c.JSON(422, gin.H{
			"error": ERR_LENGTH,
			"message": fmt.Sprintf(
				"must have len(x) == len(y) && maxDeg < len(x, y) <= %v, got len(x) = %v and len(y) = %v", MAX_ELEMENTS, n, m,
			),
		})
		return
	}

	coef, err := polyfit.PolynomialRegression(*X.X, *X.Y, int(*X.MaxDeg))
	if err != nil {
		c.JSON(500, gin.H{"error": ERR_GENERAL, "message": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(200, gin.H{"coef": coef})
}
