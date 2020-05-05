package premium

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func AddRoutes(r *gin.Engine, apiBase string, jwtMiddleware *jwt.GinJWTMiddleware) {
	group := r.Group(apiBase)
	group.Use(jwtMiddleware.MiddlewareFunc())

	r.POST("/dcf-valuation", DcfValuation)
}
