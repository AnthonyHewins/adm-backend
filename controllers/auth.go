package controllers

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/AnthonyHewins/adm-backend/models"
)

var identityKey = "id"

type login struct {
	Email    string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func genAuthMiddleware(privkey, pubkey string) *jwt.GinJWTMiddleware {
	log.Printf("Bootstrapping JWT encryption with privkey '%v' and pubkey '%v'\n", privkey, pubkey)

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: "/tools",
		MaxRefresh:  time.Hour,

		SigningAlgorithm: "RS256",
		PrivKeyFile:      privkey,
		PubKeyFile:       pubkey,

		IdentityKey:     identityKey,

		Unauthorized:    unauthorizedHandler,

		Authenticator: authenticate,
		Authorizator:  authorizator,
	})

	if err != nil { log.Fatalln(err) }

	return authMiddleware
}

func unauthorizedHandler(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error":   ERR_UNAUTHORIZED,
		"message": message,
	})
}

func authenticate(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	db, err := models.Connect()
	if err != nil { return "", err }

	user := models.User{}
	if db.Where("email = ?", loginVals.Email).First(&user).RecordNotFound() {
		return nil, jwt.ErrFailedAuthentication
	}

	if user.ConfirmedAt == nil {
		return nil, fmt.Errorf("you have not confirmed your email yet; please confirm it to log in")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(loginVals.Password),
	)

	switch err {
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, jwt.ErrFailedAuthentication
	case nil:
		return gin.H{identityKey: user.ID}, nil
	default:
		return nil, err
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
