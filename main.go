package main

import (
	"github.com/gin-gonic/gin"
	"github.com/AnthonyHewins/adm/endpoints"
)

const (
	POLYREG             = "/api/tools/poly-reg"
	FEATURE_ENGINEERING = "/api/tools/feature-engineering"
)

func routerSetup() *gin.Engine {
	r := gin.Default()
	r.POST(POLYREG, tools.PolynomialRegression)
	r.POST(FEATURE_ENGINEERING, tools.FeatureEngineering)
	return r
}

func main() {
	r := routerSetup()
	r.Run()
}
