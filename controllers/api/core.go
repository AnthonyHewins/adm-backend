package api

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/AnthonyHewins/adm-backend/models"
)

const (
	ErrFormIncorrect = "form-incorrect"
	ErrNotFound = "not-found"
	ErrServer = "server"
)

type Payload interface {
	ToPayload() gin.H
}

type Affirmative struct {
	Msg string
}

func (a Affirmative) ToPayload() gin.H {
	return gin.H{ "message": a.Msg }
}

type Error struct {
	Http int
	Code string
	Msg  string
}

func (e Error) ToPayload() gin.H {
	return gin.H{
		"error": e.Code,
		"message": e.Msg,
	}
}

type ApiEndpoint func(c *gin.Context) (Payload, *Error)
type FormLambda  func() (Payload, *Error)
type DBLambda    func(db *gorm.DB) (Payload, *Error)

func Endpoint(fn ApiEndpoint) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := fn(c)

		if err != nil {
			c.JSON(err.Http, err.ToPayload())
		} else if payload != nil {
			c.JSON(200, payload.ToPayload())
		}
	}
}

func RequireBind(c *gin.Context, structToBind interface{}, fn FormLambda) (Payload, *Error) {
	if err := c.BindJSON(structToBind); err != nil {
		return nil, &Error{Http: 400, Code: ErrFormIncorrect, Msg: err.Error()}
	}

	return fn()
}

func RequireDB(c *gin.Context, fn DBLambda) (Payload, *Error) {
	db, err := models.Connect()

	if err != nil {
		return nil, &Error{Http: 500, Code: ErrServer, Msg: err.Error()}
	}

	defer db.Close()

	return fn(db)
}

func RequireBindAndDB(c *gin.Context, structToBind interface{}, fn DBLambda) (Payload, *Error) {
	// Recursion doesn't appear to work here...the compiler won't allow it, so this
	// suboptimal call first then call the next is what I had to settle with
	return RequireBind(c, structToBind, func() (Payload, *Error) {
		return RequireDB(c, fn)
	})
}
