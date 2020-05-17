package free

import (
	"fmt"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/feature-scaling"
	"github.com/gin-gonic/gin"
)

const (
	errTooManyElements = "too many elements; based on the first row, determined matrix to be %v x %v > max of %v"
)

type featureEngineeringReq struct {
	X    *[][]float64 `form:"x"    json:"x"    binding:"required"`
	Mode *string      `form:"mode" json:"mode" binding:"required"`
}

func (f *featureEngineeringReq) ToPayload() gin.H {
	return gin.H{"x": f.X}
}

func FeatureEngineering(c *gin.Context) (api.Payload, *api.Error) {
	var X featureEngineeringReq

	return api.RequireBind(c, &X, func() (api.Payload, *api.Error) {
		return featureEngineering(&X)
	})
}

func featureEngineering(x *featureEngineeringReq) (api.Payload, *api.Error) {
	a := *x.X
	n := len(a)

	if n <= 1 {
		return x, nil
	}

	m := len(a[0])

	if m*n > maxElements {
		return nil, &api.Error{Http: 422, Code: ErrLength, Msg: fmt.Sprintf(errTooManyElements, n, m, maxElements)}
	}

	for i := 0; i < n; i++ {
		current := len(a[i])
		if current != m {
			msg := fmt.Sprintf("need a rectangular matrix, but row %v had length %v, expected length %v", i, current, m)
			return nil, &api.Error{Http: 422, Code: ErrLength, Msg: msg}
		}
	}

	switch *x.Mode {
	case "zscore":
		verticalMap(n, m, *x.X, fe.ZScore)
	case "mean-normalization":
		verticalMap(n, m, *x.X, fe.MeanNormalization)
	default:
		return nil, &api.Error{Http: 422, Code: ErrCmd, Msg: fmt.Sprintf("don't understand mode %v", *x.Mode)}
	}

	return x, nil
}

func verticalMap(n, m int, x [][]float64, fn func([]float64)) {
	for i := 0; i < m; i++ {

		buf := make([]float64, n)
		for j := 0; j < n; j++ {
			buf[j] = x[j][i]
		}

		fn(buf)
		for j := 0; j < n; j++ {
			x[j][i] = buf[j]
		}
	}
}
