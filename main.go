package main

import (
	"github.com/gin-gonic/gin"
	"github.com/AnthonyHewins/adm-backend/controllers"
)

const (
	POLYREG             = "/api/controllers/poly-reg"
	FEATURE_ENGINEERING = "/api/controllers/feature-engineering"
)

func routerSetup() *gin.Engine {
	r := gin.Default()

	r.POST(POLYREG, controllers.PolynomialRegression)
	r.POST(FEATURE_ENGINEERING, controllers.FeatureEngineering)

	return r
}

func main() {
	r := routerSetup()
	r.Run()
}
