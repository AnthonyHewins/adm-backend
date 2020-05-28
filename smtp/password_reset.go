package smtp

import (
	"github.com/matcornic/hermes/v2"
)

func PasswordReset(email, token string) error {
	htmlMarkup := hermes.Email{
		Body: hermes.Body{
			Name:   email,
			Intros: []string{"Reset your password"},
			Actions: []hermes.Action{
				{
					Instructions: "Copy this token and paste it into the password reset form.",
					InviteCode: token,
				},
			},
			Outros: []string{
				// TODO something else here
				"Need help, or have questions? Just reply to this email.",
			},
		},
	}

	emailBody, err := smtpMasterConfig.theme.GenerateHTML(htmlMarkup)
	if err != nil {
		return err
	}

	return sendEmail(email, "Reset your password", emailBody)
}
