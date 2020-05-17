package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/AnthonyHewins/adm-backend/controllers"
	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/AnthonyHewins/adm-backend/smtp"
)

const defaultConfigFile = "./server-config.yml"

type config struct {
	AppName string `yaml:"appName"`
	BaseUrl string `yaml:"baseUrl"`
	Privkey string `yaml:"privkey"`
	Pubkey  string `yaml:"pubkey"`
	ApiBase string `yaml:"apiBase"`

	Smtp smtp.Smtp `yaml:"smtp"`
	DB   models.DB `yaml:"db"`
}

func readConfig(file *string) config {
	fptr, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalln(err)
	}

	var c config
	if err := yaml.Unmarshal(fptr, &c); err != nil {
		log.Fatalln(err)
	}

	return c
}

func main() {
	configFile := flag.String("config", defaultConfigFile, fmt.Sprintf("Configuration file to use. Default: %v", defaultConfigFile))
	flag.Parse()

	//=======================================================================
	// 1. Read config
	//=======================================================================
	log.Println("Reading config...")
	c := readConfig(configFile)
	log.Println("Done reading.")
	log.Printf("App name is %v, and baseUrl is %v\n", c.AppName, c.BaseUrl)

	//=======================================================================
	// 2. Email setup
	//=======================================================================
	log.Println("Setting up email settings...")
	smtp.EmailSetup(&c.Smtp, c.AppName, c.BaseUrl)
	log.Println("Email set up.")
	log.Printf("Sending emails from user '%v', domain '%v:%v'\n", c.Smtp.Email, c.Smtp.Domain, c.Smtp.Port)

	//=======================================================================
	// 3. Set up DB connection and test it
	//=======================================================================
	log.Println("Setting up DB...")
	models.DBSetup(&c.DB)
	log.Printf("DB config set up. Spinning up DB %v on %v:%v with user %v\n", c.DB.Name, c.DB.Host, c.DB.Port, c.DB.User)

	log.Println("Quickly testing connection...")
	db, err := models.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	db.Close()
	log.Println("Connection verified.")

	//=======================================================================
	// 4. Bind server, and finally run it
	//=======================================================================
	log.Println("Binding routes...")

	r := controllers.Router(c.ApiBase, c.Privkey, c.Pubkey)

	log.Println("Server starting.")

	r.Run()
}
