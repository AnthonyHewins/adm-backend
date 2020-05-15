package auth

import (
	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type pwResetReqConfirm struct {
	Token    string `json:"token"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func ConfirmPwReset(c *gin.Context) (api.Payload, *api.Error) {
	var resetForm pwResetReqConfirm
	return api.RequireBindAndDB(c, &resetForm, func(db *gorm.DB) (api.Payload, *api.Error) {
		return confirmPwReset(db, &resetForm)
	})
}

func confirmPwReset(db *gorm.DB, resetForm *pwResetReqConfirm) (api.Payload, *api.Error) {
	upr := models.UserPasswordReset{Token: resetForm.Token}
	err := upr.ResetPassword(db, resetForm.Password)

	switch err {
	case models.TokenTimeout:
		return nil, &api.Error{Http: 403, Code: ErrLate, Msg: err.Error()}
	case gorm.ErrRecordNotFound:
		return nil, &api.Error{Http: 401, Code: ErrUnauthorized, Msg: "unauthorized"}
	case models.PasswordTooSimple:
		return nil, &api.Error{Http: 422, Code: ErrWeakPassword, Msg: err.Error()}
	case nil:
		return &api.Affirmative{Msg: "password reset successfully"}, nil
	default:
		return nil, &api.Error{Http: 500, Code: api.ErrServer, Msg: err.Error()}
	}
}
