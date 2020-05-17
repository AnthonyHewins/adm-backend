package auth

/*
import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/adm-backend/models"
	"gopkg.in/yaml.v2"
)

func Test_passwordReset(t *testing.T) {
	var config *models.DB
	f, _ := os.Open("../../testdata/dbconfig.yml")
	buf, _ := ioutil.ReadAll(f)
	yaml.Unmarshal(buf, config)
	models.DBSetup(config)

	db, _ := models.Connect()

	notConfirmed := models.User{Email: fmt.Sprintf("r42r%v@sdmjoif.co", time.Now().UnixNano()), Password: "dsouifj"}
	notConfirmed.Create(db)

	confirmed := models.User{Email: fmt.Sprintf("r92r%v@sdmjoif.co", time.Now().UnixNano()), Password: "dsouifj"}
	confirmed.Create(db)

	tests := []struct {
		name  string
		args  pwResetReq
		want  api.Payload
		want1 *api.Error
	}{
		{
			"Fake RecordNotFound with affirmative",
			pwResetReq{Email: "suhf"},
			api.Affirmative{Msg: "if the account for suhf exists, an email has been sent to reset the password"},
			nil,
		},
		{
			"Not confirmed email fakes affirmative (this should be different)",
			pwResetReq{Email: notConfirmed.Email},
			api.Affirmative{Msg: fmt.Sprintf("if the account for %v exists, an email has been sent to reset the password", notConfirmed.Email)},
			nil,
		},
		{
			"Fake RecordNotFound with affirmative",
			pwResetReq{Email: "suhf"},
			api.Affirmative{Msg: "if the account for suhf exists, an email has been sent to reset the password"},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := passwordReset(db, &tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("passwordReset() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("passwordReset() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
*/
