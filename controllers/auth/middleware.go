package auth

import (
	"log"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/AnthonyHewins/adm-backend/models"
)

var identityKey = "id"

func GenAuthMiddleware(privkey, pubkey string) *jwt.GinJWTMiddleware {
	log.Printf("Bootstrapping JWT encryption with privkey '%v' and pubkey '%v'\n", privkey, pubkey)

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "/tools",
		MaxRefresh: time.Hour,

		SigningAlgorithm: "RS256",
		PrivKeyFile:      privkey,
		PubKeyFile:       pubkey,

		IdentityKey: identityKey,

		Unauthorized: unauthorizedHandler,

		Authenticator: authenticate,
		Authorizator:  authorizator,
	})

	if err != nil {
		log.Fatalln(err)
	}

	return authMiddleware
}

func unauthorizedHandler(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error":   ErrUnauthorized,
		"message": message,
	})
}

func authenticate(c *gin.Context) (interface{}, error) {
	var loginVals credentials
	if err := c.ShouldBind(&loginVals); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	db, err := models.Connect()
	if err != nil {
		return nil, err
	}

	user := models.User{Email: loginVals.Email, Password: loginVals.Password}
	err = user.Authenticate(db)

	switch err {
	case models.EmailNotConfirmed:
		return nil, err
	case nil:
		return gin.H{identityKey: user.ID}, nil
	default:
		return nil, jwt.ErrFailedAuthentication
	}
}

func authorizator(data interface{}, c *gin.Context) bool {
	// TODO when the route permissions start sharing logic
	// if v, ok := data.(*models.User); ok && v.Email == "admin" {
	// }

	return true
}

func logoutHandler(c *gin.Context, code int) {
	c.JSON(200, gin.H{
		"message": "successfully logged out",
	})
}
