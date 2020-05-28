package smtp

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
)

var smtpMasterConfig masterConfig

type masterConfig struct {
	smtpSettings Smtp
	baseUrl      string

	theme hermes.Hermes
	send  bool
}

type Smtp struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Domain   string `yaml:"domain"`
	Port     int    `yaml:"port"`

	ConfirmationUrl string `yaml:"confirmationUrl"`

	Send bool `yaml:"send"`
}

func EmailSetup(smtpSettings *Smtp, appName, baseUrl string) {
	smtpMasterConfig = masterConfig{
		smtpSettings: *smtpSettings,
		baseUrl:      baseUrl,
		theme: hermes.Hermes{
			Product: hermes.Product{
				Name: appName,
				Link: baseUrl,
				Logo: fmt.Sprintf("%v/favicon.ico", baseUrl),
			},
		},
	}
}

func sendEmail(to, subject, body string) error {
	if !smtpMasterConfig.smtpSettings.Send {
		log.Println("Sending is OFF; faked an email send.")
		return nil
	}

	m := gomail.NewMessage()

	m.SetHeader("From", smtpMasterConfig.smtpSettings.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		smtpMasterConfig.smtpSettings.Domain,
		smtpMasterConfig.smtpSettings.Port,
		smtpMasterConfig.smtpSettings.Email,
		smtpMasterConfig.smtpSettings.Password,
	)
	d.TLSConfig = &tls.Config{ InsecureSkipVerify: true	}

	return d.DialAndSend(m)
}
