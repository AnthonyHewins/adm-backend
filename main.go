package main

import (
	"github.com/gin-gonic/gin"
	"github.com/AnthonyHewins/adm/endpoints"
)

func main() {
	r := gin.Default()
	r.POST("/tools/poly-reg", tools.PolynomialRegression)
	r.Run() // listen and serve on 0.0.0.0:8080
}

