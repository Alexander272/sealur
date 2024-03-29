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
		Postgres    PostgresConfig
		Http        HttpConfig
		Api         ApiConfig
		Services    ServicesConfig
	}

	PostgresConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string
		DbName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
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

	ServicesConfig struct {
		EmailService ServiceConfig
		FileService  ServiceConfig
		UserService  ServiceConfig
	}

	ServiceConfig struct {
		Url          string
		AuthName     string
		AuthPassword string
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
	if err := viper.UnmarshalKey("postgres", &conf.Postgres); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("http", &conf.Http); err != nil {
		return err
	}

	return nil
}

func setFromEnv(conf *Config) error {
	if err := envconfig.Process("postgres", &conf.Postgres); err != nil {
		return err
	}
	// if err := envconfig.Process("http", &conf.Http); err != nil {
	// 	return err
	// }
	if err := envconfig.Process("api", &conf.Api); err != nil {
		return err
	}

	conf.Http.Host = os.Getenv("PRO_HOST")
	conf.Http.Port = os.Getenv("PRO_PORT")
	if len(conf.Http.Port) == 0 {
		conf.Http.Port = "9001"
	}
	conf.Postgres.Password = os.Getenv("PRO_DB_PASSWORD")
	conf.Environment = os.Getenv("APP_ENV")

	conf.Services.EmailService.AuthName = os.Getenv("API_NAME")
	conf.Services.EmailService.AuthPassword = os.Getenv("API_PASSWORD")
	conf.Services.FileService.AuthName = os.Getenv("API_NAME")
	conf.Services.FileService.AuthPassword = os.Getenv("API_PASSWORD")
	conf.Services.UserService.AuthName = os.Getenv("API_NAME")
	conf.Services.UserService.AuthPassword = os.Getenv("API_PASSWORD")

	conf.Services.EmailService.Url = os.Getenv("EMAIL_HOST") + ":" + os.Getenv("EMAIL_PORT")
	conf.Services.FileService.Url = os.Getenv("FILE_HOST") + ":" + os.Getenv("FILE_PORT")
	conf.Services.UserService.Url = os.Getenv("USER_HOST") + ":" + os.Getenv("USER_PORT")

	return nil
}
