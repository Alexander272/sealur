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
		Redis       RedisConfig
		Auth        AuthConfig
		Http        HttpConfig
		Limiter     LimiterConfig

		Services ServicesConfig
	}

	RedisConfig struct {
		Host     string `mapstructure:"Host"`
		Port     string `mapstructure:"Port"`
		DB       int    `mapstructure:"DB"`
		Password string
	}

	AuthConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		Secure          bool          `mapstructure:"secure"`
		Domain          string        `mapstructure:"domain"`
		Key             string
	}

	HttpConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	LimiterConfig struct {
		RPS   int           `mapstructure:"rps"`
		Burst int           `mapstructure:"burst"`
		TTL   time.Duration `mapstructure:"ttl"`
	}

	ServicesConfig struct {
		ProService  ProConfig
		FileService ProConfig
		UserService ProConfig
	}

	ProConfig struct {
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
	if err := viper.UnmarshalKey("redis", &conf.Redis); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("http", &conf.Http); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("limiter", &conf.Limiter); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &conf.Auth); err != nil {
		return err
	}

	return nil
}

func setFromEnv(conf *Config) error {
	if err := envconfig.Process("http", &conf.Http); err != nil {
		return err
	}
	if err := envconfig.Process("jwt", &conf.Auth); err != nil {
		return err
	}
	if err := envconfig.Process("redis", &conf.Redis); err != nil {
		return err
	}

	conf.Environment = os.Getenv("APP_ENV")
	conf.Services.ProService.AuthName = os.Getenv("API_NAME")
	conf.Services.ProService.AuthPassword = os.Getenv("API_PASSWORD")
	conf.Services.ProService.Url = os.Getenv("PRO_HOST") + ":" + os.Getenv("PRO_PORT")

	conf.Services.FileService.AuthName = os.Getenv("API_NAME")
	conf.Services.FileService.AuthPassword = os.Getenv("API_PASSWORD")
	conf.Services.FileService.Url = os.Getenv("FILE_HOST") + ":" + os.Getenv("FILE_PORT")

	conf.Services.UserService.AuthName = os.Getenv("API_NAME")
	conf.Services.UserService.AuthPassword = os.Getenv("API_PASSWORD")
	conf.Services.UserService.Url = os.Getenv("USER_HOST") + ":" + os.Getenv("USER_PORT")

	return nil
}
