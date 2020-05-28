package smtp

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
)

func AccountConfirmation(email, token string) error {
	return transactional(
		email,
		token,
		"Welcome to Artifex de Machina!",
		"To confirm your account all you need to do is click the link on the button below.",
	)
}

func TokenRefresh(email, token string) error {
	return transactional(
		email,
		token,
		"Confirm your email",
		"We sent another email to you because your email was validated too late or because you requested another one be sent. To confirm your account, click the button below.",
	)
}

func transactional(email, token, header, instructions string) error {
	link := fmt.Sprintf(
		// Note: first two args are paths so they end in /
		// baseUrl/ + confirmUrl/ + "/" + token
		"%v%v/%v",
		smtpMasterConfig.baseUrl,
		smtpMasterConfig.smtpSettings.ConfirmationUrl,
		token,
	)

	htmlMarkup := hermes.Email{
		Body: hermes.Body{
			Name:   email,
			Intros: []string{header},
			Actions: []hermes.Action{
				{
					Instructions: instructions,
					Button: hermes.Button{
						//Color: "#22BC66",
						Text: "Continue",
						Link: link,
					},
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

	return sendEmail(email, header, emailBody)
}
