package config

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Environment string
		Http        HttpConfig
		Api         ApiConfig
		Email       EmailConfig
		Recipients  RecipientsConfig
	}

	HttpConfig struct {
		ServiceName        string        `mapstructure:"serviceName"`
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	ApiConfig struct {
		Name     string
		Password string
	}

	EmailConfig struct {
		Sender   string
		Password string
		Host     string
		Port     int
	}

	RecipientsConfig struct {
		Interview string
		Confirm   string
		Order     string
		Test      string
		Support   string
		Link      string
	}
)

func Init(configDir string) (*Config, error) {
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}

	var conf Config
	if err := unmarhal(&conf); err != nil {
		return nil, err
	}
	if err := setFromEnv(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	return viper.MergeInConfig()
}

func unmarhal(conf *Config) error {
	if err := viper.UnmarshalKey("http", &conf.Http); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("recipients", &conf.Recipients); err != nil {
		return err
	}

	return nil
}

func setFromEnv(conf *Config) error {
	if err := envconfig.Process("api", &conf.Api); err != nil {
		return err
	}
	if err := envconfig.Process("smtp", &conf.Email); err != nil {
		return err
	}

	conf.Http.Host = os.Getenv("EMAIL_HOST")
	conf.Http.Port = os.Getenv("EMAIL_PORT")
	if len(conf.Http.Port) == 0 {
		conf.Http.Port = "12001"
	}
	conf.Environment = os.Getenv("APP_ENV")

	return nil
}
