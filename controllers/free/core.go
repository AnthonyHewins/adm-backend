package free

import (
	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/gin-gonic/gin"
)

const (
	maxElements = 100

	ErrCmd    = "cmd"
	ErrLength = "length"
	ErrDegree = "deg"
)

func AddRoutes(r *gin.Engine, apiBase string) {
	group := r.Group(apiBase)
	group.POST("/poly-reg",            api.Endpoint(PolynomialRegression))
	group.POST("/feature-engineering", api.Endpoint(FeatureEngineering))
}
