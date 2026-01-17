package config

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

// pake -ldflags di perintah build untuk set variable dibawah
//
// contoh: -ldflags="-X github.com/crazydw4rf/oil-bank-backend/internal/config.APP_VERSION=1.0.4-beta"
var (
	APP_ENV          = "development"
	APP_VERSION      = "0.0.1-alpha"
	API_DOCS_ENABLED = true
)

const (
	BASE_API_HTTP_PATH            = "/v1"
	REFRESH_TOKEN_COOKIE_NAME     = "refreshToken"
	ACCESS_TOKEN_COOKIE_NAME      = "token"
	ACCESS_TOKEN_HEADER_NAME      = "Authorization"
	ACCESS_TOKEN_EXPIRATION_TIME  = time.Minute * 15
	REFRESH_TOKEN_EXPIRATION_TIME = (time.Hour * 24) * 30
)

type Config struct {
	APP_PORT                 int    `mapstructure:"APP_PORT"`
	APP_HOST                 string `mapstructure:"APP_HOST"`
	DATABASE_URL             string `mapstructure:"DATABASE_URL"`
	JWT_ACCESS_TOKEN_SECRET  string `mapstructure:"JWT_ACCESS_TOKEN_SECRET"`
	JWT_REFRESH_TOKEN_SECRET string `mapstructure:"JWT_REFRESH_TOKEN_SECRET"`
}

func InitConfig() (*Config, error) {
	cfg := new(Config)
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigType("env")
	v.SetConfigName(".env")
	v.AutomaticEnv()

	bindEnvStruct(v, cfg)

	err := v.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning reading config file: %#v\n", err)
	}

	// TODO: struct validation?
	err = v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func bindEnvStruct(v *viper.Viper, s any) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()

	for i := range typ.NumField() {
		field := typ.Field(i)
		tagValue := field.Tag.Get("mapstructure")
		if tagValue != "" {
			v.BindEnv(field.Name, tagValue)
		}
	}
}
