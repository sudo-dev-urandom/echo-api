package core

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kelseyhightower/envconfig"
)

var (
	App *Application
)

type (
	Application struct {
		Name    string   `json:"name"`
		Port    string   `json:"port"`
		Version string   `json:"version"`
		Config  Config   `json:"app_config"`
		DB      *gorm.DB `json:"db"`
	}

	Config struct {
		Port                         string `envconfig:"APPPORT"`
		DB_HOST                      string `envconfig:"DB_HOST"`
		DB_USER                      string `envconfig:"DB_USER"`
		DB_PASS                      string `envconfig:"DB_PASS"`
		DB_NAME                      string `envconfig:"DB_NAME"`
		DB_PORT                      string `envconfig:"DB_PORT"`
		DB_LOG                       int    `envconfig:"DB_LOG"`
		JWT_SECRET                   string `envconfig:"JWT_SECRET"`
		GOOGLE_CLIENT_ID             string `envconfig:"GOOGLE_CLIENT_ID"`
		GOOGLE_CLIENT_SECRET         string `envconfig:"GOOGLE_CLIENT_SECRET"`
		GOOGLE_CLIENT_REDIRECT_URL   string `envconfig:"GOOGLE_CLIENT_REDIRECT_URL"`
		FACEBOOK_CLIENT_ID           string `envconfig:"FACEBOOK_CLIENT_ID"`
		FACEBOOK_CLIENT_SECRET       string `envconfig:"FACEBOOK_CLIENT_SECRET"`
		FACEBOOK_CLIENT_REDIRECT_URL string `envconfig:"FACEBOOK_CLIENT_REDIRECT_URL"`
	}
)

func init() {
	var err error
	App = &Application{}

	if err = App.LoadConfigs(); err != nil {
		log.Printf("Load config error conf: %v", err)
	}

	if err = App.DatabaseInit(); err != nil {
		log.Printf("Load config error db: %v", err)
	}
}

func (x *Application) LoadConfigs() error {

	err := envconfig.Process("myapp", &x.Config)
	x.Name = "Example App"
	x.Version = os.Getenv("APPVER")
	x.Port = x.Config.Port

	return err
}

func (x *Application) DatabaseInit() error {
	config := x.Config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", config.DB_USER, config.DB_PASS, config.DB_HOST, config.DB_PORT, config.DB_NAME)

	db, err := gorm.Open("mysql", dsn)
	db.LogMode(config.DB_LOG == 1)
	x.DB = db

	return err
}

func (x *Application) Close() (err error) {
	if err = x.DB.Close(); err != nil {
		return err
	}

	return nil
}
