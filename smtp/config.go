package smtp

import (
	"fmt"
	"github.com/matcornic/hermes/v2"
)

var smtpMasterConfig masterConfig

type masterConfig struct {
	smtpSettings         Smtp
	baseUrl              string
	confirmationEndpoint string

	theme                hermes.Hermes
	send                 bool
}

type Smtp struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Domain   string `yaml:"domain"`
	Port     int    `yaml:"port"`

	Send     bool   `yaml:"send"`
}

func EmailSetup(smtpSettings *Smtp, appName, baseUrl, confirmationEndpoint string) {
	smtpMasterConfig = masterConfig{
		smtpSettings: *smtpSettings,
		baseUrl: baseUrl,
		confirmationEndpoint: confirmationEndpoint,
		theme: hermes.Hermes{
			Product: hermes.Product{
				Name: appName,
				Link: baseUrl,
				Logo: fmt.Sprintf("%v/favicon.ico", baseUrl),
			},
		},
	}
}
