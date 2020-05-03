package smtp

import (
	"fmt"
	"log"

	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
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

func PasswordReset(email, token string) error {
	return transactional(
		email,
		token,
		"Reset your password",
		"You requested to reset your password. To reset it, click the link below.",
	)
}

func transactional(email, token, header, instructions string) error {
	if !smtpMasterConfig.smtpSettings.Send {
		log.Println("Sending is OFF; faked an email send.")
		return nil
	}

	link := fmt.Sprintf(
		"%v%v?token=%v",
		smtpMasterConfig.baseUrl,
		smtpMasterConfig.smtpSettings.ConfirmationUrl,
		token,
	)

	htmlMarkup := hermes.Email{
		Body: hermes.Body{
			Name: email,
			Intros: []string{ header },
			Actions: []hermes.Action{
				{
					Instructions: instructions,
					Button: hermes.Button{
						//Color: "#22BC66",
						Text:  "Continue",
						Link:  link,
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
	if err != nil { return err }

	m := gomail.NewMessage()

	m.SetHeader("From", smtpMasterConfig.smtpSettings.Email)
	m.SetHeader("To", email)
	m.SetHeader("Subject", header)
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(
		smtpMasterConfig.smtpSettings.Domain,
		smtpMasterConfig.smtpSettings.Port,
		smtpMasterConfig.smtpSettings.Email,
		smtpMasterConfig.smtpSettings.Password,
	)

	return d.DialAndSend(m)
}
