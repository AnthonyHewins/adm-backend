package auth

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/adm-backend/models"
)

var (
	uniquenessViolation = regexp.MustCompile("duplicate key value violates unique constraint")
)

func Register(c *gin.Context) (api.Payload, *api.Error) {
	var form credentials

	return api.RequireBindAndDB(c, &form, func(db *gorm.DB) (api.Payload, *api.Error) {
		return register(db, &form)
	})
}

func register(db *gorm.DB, req *credentials) (api.Payload, *api.Error) {
	u := req.toUser()
	err := u.Create(db)

	switch err {
	case models.InvalidEmail:
		return nil, &api.Error{Http: 422, Code: ErrEmail, Msg: err.Error()}
	case models.PasswordTooSimple:
		return nil, &api.Error{Http: 422, Code: ErrWeakPassword, Msg: err.Error()}
	case nil:
		return &api.Affirmative{Msg: "email sent; please confirm it with the email that was sent to your address"}, nil
	default:
		errstring := err.Error()

		if uniquenessViolation.MatchString(errstring) {
			return nil, &api.Error{Http: 422, Code: ErrEmailTaken, Msg: "email is taken, please use a different one"}
		} else {
			return nil, &api.Error{Http: 500, Code: api.ErrServer, Msg: errstring}
		}
	}
}
