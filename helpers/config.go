package helpers

import (
	"os"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App       App
		DB        DB
		JWTSecret string `mapstructure:"JWT_SECRET"`
		PhotoDir  string `mapstructure:"PHOTO_DIR"`
	}
	App struct {
		Port uint `mapstructure:"APP_PORT"`
	}
	DB struct {
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
		Name     string `mapstructure:"DB_NAME"`
		Host     string `mapstructure:"DB_HOST"`
		Port     uint   `mapstructure:"DB_PORT"`
	}
)

func LoadConfig(configFile string) (Config, error) {
	var (
		app  App
		db   DB
		conf Config
	)
	_, err := os.Stat(configFile)
	if err != nil {
		return conf, err
	}

	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigFile(configFile)

	if err := v.ReadInConfig(); err != nil {
		return conf, err
	}

	if err := v.Unmarshal(&db); err != nil {
		return conf, err
	}

	if err := v.Unmarshal(&app); err != nil {
		return conf, err
	}

	if err := v.Unmarshal(&conf); err != nil {
		return conf, err
	}

	conf.DB = db
	conf.App = app

	os.Setenv("JWT_SECRET", conf.JWTSecret)
	os.Setenv("PHOTO_DIR", conf.PhotoDir)

	return conf, nil
}
