package free

import (
	"fmt"

	"github.com/AnthonyHewins/adm-backend/api"
	"github.com/AnthonyHewins/feature-scaling"
	"github.com/gin-gonic/gin"
)

type featureEngineeringData struct {
	X    *[][]float64 `form:"x"    json:"x"    binding:"required"`
	Mode *string      `form:"mode" json:"mode" binding:"required"`
}

func FeatureEngineering(c *gin.Context) {
	var X featureEngineeringData

	api.RequireBind(c, &X, func() {
		a := *X.X
		n := len(a)

		if n <= 1 {
			c.JSON(204, nil)
			return
		}

		m := len(a[0])

		if m * n > maxElements {
			msg := fmt.Sprintf("too many elements; based on the first row, determined matrix to be %v x %v > max of %v", n, m, maxElements)
			api.Error(c, 422, msg)
			return
		}

		for i := 0; i < n; i++ {
			current := len(a[i])
			if current != m {
				msg := fmt.Sprintf("need a rectangular matrix, but row %v had length %v, expected length %v", i, current, m)
				api.Error(c, 422, msg)
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
			api.Error(c, 422, fmt.Sprintf("don't understand mode %v", *X.Mode))
		}
	})

}

func verticalMap(n, m int, x [][]float64, fn func([]float64)) {
	for i := 0; i < m; i++ {

		buf := make([]float64, n)
		for j := 0; j < n; j++ { buf[j]  = x[j][i] }

		fn(buf)
		for j := 0; j < n; j++ { x[j][i] = buf[j]  }
	}
}
