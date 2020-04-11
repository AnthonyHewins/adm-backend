package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/AnthonyHewins/feature-scaling"
)

type FeatureEngineeringData struct {
	X    *[][]float64 `form:"x"    json:"x"    binding:"required"`
	Mode *string      `form:"mode" json:"mode" binding:"required"`
}

func FeatureEngineering(c *gin.Context) {
	var X FeatureEngineeringData

	if err := c.BindJSON(&X); err != nil {
		c.JSON(400, gin.H{
			"error": ERR_PARAM,
			"message": err.Error(),
		})
		return
	}

	a := *X.X
	n := len(a)

	if n <= 1 {
		c.JSON(204, nil)
		return
	}

	m := len(a[0])

	if m * n > MAX_ELEMENTS {
		c.JSON(422, gin.H{
			"error": ERR_LENGTH,
			"message": fmt.Sprintf("too many elements; based on the first row, determined matrix to be %v x %v > max of %v", n, m, MAX_ELEMENTS),
		})
		return
	}

	for i := 0; i < n; i++ {
		current := len(a[i])
		if current != m {
			c.JSON(422, gin.H{
				"error": ERR_LENGTH_MISMATCH,
				"message": fmt.Sprintf("need a rectangular matrix, but row %v had length %v, expected length %v", i, current, m),
			})
			return
		}
	}

	switch *X.Mode {
	case "zscore":
		verticalMap(n, m, *X.X, fe.ZScore)
		c.JSON(200, gin.H{
			"x": X.X,
		})
	case "mean-normalization":
		verticalMap(n, m, *X.X, fe.MeanNormalization)
		c.JSON(200, gin.H{
			"x": X.X,
		})
	default:
		c.JSON(422, gin.H{
			"error": ERR_CMD,
			"message": fmt.Sprintf("don't understand mode %v", *X.Mode),
		})
	}
}

func verticalMap(n, m int, x [][]float64, fn func([]float64)) {
	for i := 0; i < m; i++ {

		buf := make([]float64, n)
		for j := 0; j < n; j++ { buf[j]  = x[j][i] }

		fn(buf)
		for j := 0; j < n; j++ { x[j][i] = buf[j]  }
	}
}
