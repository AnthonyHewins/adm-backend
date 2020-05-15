package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/AnthonyHewins/adm-backend/smtp"
)

func AcctConfirmation(c *gin.Context) (api.Payload, *api.Error) {
	token := c.Query("token")

	// Don't let a botnet know this is a valid link
	if token == "" {
		c.String(404, "page not found")
		return nil, nil
	}

	return api.RequireDB(c, func(db *gorm.DB) (api.Payload, *api.Error) {
		return acctConfirmation(db, token)
	})
}

func acctConfirmation(db *gorm.DB, token string) (api.Payload, *api.Error) {
	uec := models.UserEmailConfirmation{}

	err := db.Where("token = ?", token).First(&uec).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return nil, &api.Error{Http: 404, Code: api.ErrNotFound, Msg: "the server can't find what you're looking for"}
	case nil:
		// no op
	default:
		return nil, &api.Error{Http: 500, Code: api.ErrServer, Msg: err.Error()}
	}

	switch err := uec.ConfirmEmail(db); err {
	case models.TokenTimeout:
		return sendNewConfirmationToken(db, uec.UserID)
	case nil:
		return &api.Affirmative{Msg: "email confirmed, welcome"}, nil
	default:
		return nil, &api.Error{Http: 500, Code: api.ErrServer, Msg: err.Error()}
	}
}

func sendNewConfirmationToken(db *gorm.DB, id uint64) (api.Payload, *api.Error) {
	u := models.User{ID: id}
	db.First(&u)

	token, err := u.RefreshConfirmationToken(db)
	if err != nil {
		return nil, &api.Error{Http: 500, Code: api.ErrServer, Msg: err.Error()}
	}

	var msg string
	if err = smtp.TokenRefresh(u.Email, token); err == nil {
		msg = "confirmed email too late; another email has been sent"
	} else {
		msg = fmt.Sprintf(
			"confirmed email too late; tried sending another email but failed: %v",
			err.Error(),
		)
	}

	return nil, &api.Error{Http: 403, Code: ErrLate, Msg: msg}
}
