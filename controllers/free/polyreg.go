package free

import (
	"fmt"

	"github.com/AnthonyHewins/adm-backend/api"
	"github.com/AnthonyHewins/polyfit"
	"github.com/gin-gonic/gin"
)

type PolyRegData struct{
	MaxDeg *int       `form:"maxDeg" json:"maxDeg" binding:"required"`
	X      *[]float64 `form:"x"      json:"x"      binding:"required"`
	Y      *[]float64 `form:"y"      json:"y"      binding:"required"`
}

func PolynomialRegression(c *gin.Context) {
	var X PolyRegData
	api.RequireBind(c, &X, func() {
		if maxDegree < *X.MaxDeg || 0 >= *X.MaxDeg {
			c.JSON(422, api.ToError(errDegree, fmt.Sprintf("maxDeg must satisfy 0 <= maxDeg <= %v", maxDegree)))
			return
		}

		n := len(*X.X)
		m := len(*X.Y)

		if (n != m || n > maxElements || n <= *X.MaxDeg) {
			msg := fmt.Sprintf("must have len(x) == len(y) && maxDeg < len(x, y) <= %v, got len(x) = %v and len(y) = %v", maxElements, n, m)
			c.JSON(422, api.ToError(errLength, msg))
			return
		}

		coef, err := polyfit.PolynomialRegression(*X.X, *X.Y, int(*X.MaxDeg))
		if err != nil {
			api.Error(c, 500, err.Error())
			return
		}

		c.JSON(200, gin.H{"coef": coef})
	})
}
