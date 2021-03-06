package auth

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/adm-backend/models"
)

func TestConfirmPwReset(t *testing.T) {
	db := dbInstance()

	// Build the state needed in the DB
	pw := "sndfojus"
	u := models.User{Email: fmt.Sprintf("r22rasdd%v@sdmjoif.co", time.Now().UnixNano()), Password: pw}
	u.Create(db)

	upr := models.UserPasswordReset{UserID: u.ID}
	err := upr.CreateResetPasswordToken(db)
	if err != nil { t.Fatalf("Failed creation of UserPasswordReset: %v", err) }

	// Build the state needed in the DB
	timeoutPw := "sndfojus"
	timedout := models.User{Email: fmt.Sprintf("r123122dr%v@sdmjoif.co", time.Now().UnixNano()), Password: timeoutPw}
	timedout.Create(db)

	timedoutUpr := models.UserPasswordReset{UserID: timedout.ID}
	timedoutUpr.CreateResetPasswordToken(db)
	db.Model(&timedoutUpr).Where("user_id = ?", timedout.ID).Update("reset_at", time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	db.First(&timedoutUpr)

	tests := []struct {
		name  string
		pwReset pwResetReqConfirm
		want  api.Payload
		want1 *api.Error
	}{
		{
			"Token doesn't exist is 401",
			pwResetReqConfirm{Token: "garbage"},
			nil,
			&api.Error{Http: 401, Code: ErrUnauthorized, Msg: "unauthorized"},
		},
		{
			"New password is too weak",
			pwResetReqConfirm{Token: upr.Token, Password: "asdff"},
			nil,
			&api.Error{Http: 422, Code: ErrWeakPassword, Msg: "password needs to be at least 6 characters"},
		},
		{
			"Token timeout",
			pwResetReqConfirm{Token: timedoutUpr.Token, Password: timeoutPw},
			nil,
			&api.Error{Http: 403, Code: ErrLate, Msg: models.TokenTimeout.Error()},
		},
		{
			"Proper reset",
			pwResetReqConfirm{Token: upr.Token, Password: pw},
			&api.Affirmative{Msg: "password reset successfully"},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := confirmPwReset(db, &tt.pwReset)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfirmPwReset() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConfirmPwReset() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	// New password works
	uec := models.UserEmailConfirmation{}
	db.Where("user_id = ?", u.ID).First(&uec)
	uec.ConfirmEmail(db)

	u.Password = pw
	if err := u.Authenticate(db); err != nil {
		t.Errorf("successful password reset is not changing the password around? Got err: %v", err)
	}
}
