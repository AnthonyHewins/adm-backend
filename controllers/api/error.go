package api

import (
	"github.com/gin-gonic/gin"
)

const (
	ErrIncorrectParams = "incorrect-params"
	ErrForbidden = "forbidden"
	ErrUnauthorized = "unauthorized"
	ErrMissing = "not-found"
	ErrInvalidParam = "param-error"
	ErrServer = "server"
)

func Error(c *gin.Context, code uint16, msg string) {
	switch code {
	case 400:
		c.JSON(400, ToError(ErrIncorrectParams, msg))
	case 401:
		c.JSON(401, ToError(ErrUnauthorized, msg))
	case 403:
		c.JSON(403, ToError(ErrForbidden, msg))
	case 404:
		c.JSON(404, ToError(ErrMissing, msg))
	case 422:
		c.JSON(422, ToError(ErrInvalidParam, msg))
	default:
		c.JSON(500, ToError(ErrServer, msg))
	}
}

func ToError(errcode, msg string) gin.H {
	return gin.H{ "error": errcode, "message": msg }
}
