package auth

import (
	"fmt"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type pwResetReq struct {
	Email string `json:"email"`
}

func PasswordReset(c *gin.Context) (api.Payload, *api.Error){
	var pwResetReqForm pwResetReq
	return api.RequireBindAndDB(c, &pwResetReqForm, func(db *gorm.DB) (api.Payload, *api.Error){
		return passwordReset(db, &pwResetReqForm)
	})
}

func passwordReset(db *gorm.DB, req *pwResetReq) (api.Payload, *api.Error) {
	user := models.User{}
	err := db.Where("email = ?", req.Email).First(&user).Error

	switch err {
	case gorm.ErrRecordNotFound:
		// Don't let an attacker know they found an email
		return fillWithOk(req.Email)
	case nil:
		// no op
	default:
		return nil, &api.Error{Http: 500, Code: api.ErrServer, Msg: err.Error()}
	}

	// Email not confirmed -> they can't reset it
	// TODO this needs to be something smarter
	if user.ConfirmedAt == nil {
		return fillWithOk(user.Email)
	}

	upr := models.UserPasswordReset{UserID: user.ID}
	if err := upr.CreateResetPasswordToken(db); err != nil {
		return nil, &api.Error{Http: 500, Code: api.ErrServer, Msg: err.Error()}
	}

	return fillWithOk(user.Email)
}

func fillWithOk(email string) (api.Payload, *api.Error) {
	msg := fmt.Sprintf("if the account for %v exists, an email has been sent to reset the password", email)
	return &api.Affirmative{Msg: msg}, nil
}
