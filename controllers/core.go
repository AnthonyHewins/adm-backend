package controllers

import (
	"log"

	"github.com/gin-gonic/gin"

	auth "github.com/AnthonyHewins/adm-backend/controllers/auth"
	free "github.com/AnthonyHewins/adm-backend/controllers/free"
	premium "github.com/AnthonyHewins/adm-backend/controllers/premium"
)

func Router(apiBase, privkey, pubkey string) *gin.Engine {
	router := gin.Default()

	log.Println("Generating middleware")
	jwtMiddleware := auth.GenAuthMiddleware(privkey, pubkey)

	log.Println("Attaching auth routes")
	auth.AddRoutes(router, apiBase, jwtMiddleware)

	log.Println("Attaching free routes")
	free.AddRoutes(router, apiBase)

	log.Println("Attaching premium routes")
	premium.AddRoutes(router, apiBase, jwtMiddleware)

	return router
}
