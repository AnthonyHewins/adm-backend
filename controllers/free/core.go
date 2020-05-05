package free

import (
	"github.com/gin-gonic/gin"
)

const (
	maxDegree   = 5
	maxElements = 100

	errCmd    = "cmd"
	errLength = "length"
	errDegree = "deg"
)

func AddRoutes(r *gin.Engine, apiBase string) {
	group := r.Group(apiBase)
	group.POST("/poly-reg",            PolynomialRegression)
	group.POST("/feature-engineering", FeatureEngineering)
}
