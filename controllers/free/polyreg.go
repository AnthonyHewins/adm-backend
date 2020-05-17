package free

import (
	"fmt"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/polyfit"
	"github.com/gin-gonic/gin"
)

const (
	maxDegree   = 5
	ErrSingular = "singular-matrix"
)

type polyRegReq struct {
	MaxDeg *int       `form:"maxDeg" json:"maxDeg" binding:"required"`
	X      *[]float64 `form:"x"      json:"x"      binding:"required"`
	Y      *[]float64 `form:"y"      json:"y"      binding:"required"`
}

type polyRegResp struct {
	Coef []float64 `json:"coef"`
}

func (p *polyRegResp) ToPayload() gin.H {
	return gin.H{"coef": p.Coef}
}

func PolynomialRegression(c *gin.Context) (api.Payload, *api.Error) {
	var X polyRegReq
	return api.RequireBind(c, &X, func() (api.Payload, *api.Error) {
		return polynomialRegression(&X)
	})
}

func polynomialRegression(X *polyRegReq) (api.Payload, *api.Error) {
	if maxDegree < *X.MaxDeg || 0 >= *X.MaxDeg {
		return nil, &api.Error{Http: 422, Code: ErrDegree, Msg: fmt.Sprintf("maxDeg must satisfy 0 <= maxDeg <= %v", maxDegree)}
	}

	n := len(*X.X)
	m := len(*X.Y)

	if n != m || n > maxElements || n <= *X.MaxDeg {
		msg := fmt.Sprintf("must have len(x) == len(y) && maxDeg < len(x, y) <= %v, got len(x) = %v and len(y) = %v", maxElements, n, m)
		return nil, &api.Error{Http: 422, Code: ErrLength, Msg: msg}
	}

	coef, err := polyfit.PolynomialRegression(*X.X, *X.Y, int(*X.MaxDeg))
	if err != nil {
		return nil, &api.Error{Http: 500, Code: ErrSingular, Msg: err.Error()}
	}

	return &polyRegResp{Coef: coef}, nil
}
