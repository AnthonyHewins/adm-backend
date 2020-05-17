package auth

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/adm-backend/models"
)

func Test_register(t *testing.T) {
	db := dbInstance()

	existingEmail := fmt.Sprintf("r2p2r%v@sdmjoif.co", time.Now().UnixNano())
	u := models.User{Email: existingEmail, Password: "sjdihfiusd"}
	u.Create(db)

	tests := []struct {
		name  string
		args  credentials
		want  api.Payload
		want1 *api.Error
	}{
		{
			"Invalid email",
			credentials{Email: "ausondfuio", Password: "asndoasd"},
			nil,
			&api.Error{Http: 422, Code: ErrEmail, Msg: "email didn't pass validation"},
		},
		{
			"weak password",
			credentials{Email: "ausondfuio", Password: "asnds"},
			nil,
			&api.Error{Http: 422, Code: ErrWeakPassword, Msg: "password needs to be at least 6 characters"},
		},
		{
			"already existing email",
			credentials{Email: existingEmail, Password: "asndoasd"},
			nil,
			&api.Error{Http: 422, Code: ErrEmailTaken, Msg: "email is taken, please use a different one"},
		},
		{
			"success",
			credentials{Email: fmt.Sprintf("ahu%v@diodsj.cod", time.Now().UnixNano()), Password: "asndoasd"},
			&api.Affirmative{Msg: "email sent; please confirm it with the email that was sent to your address"},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := register(db, &tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("register() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("register() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
