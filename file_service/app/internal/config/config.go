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
		Tcp         TcpConfig
		Api         ApiConfig
		MinIO       MinIOConfig
	}

	TcpConfig struct {
		Port string `mapstructure:"port"`
	}

	HttpConfig struct {
		Domain             string        `mapstructure:"domain"`
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

	MinIOConfig struct {
		Endpoint  string `mapstructure:"endpoint"`
		AccessKey string `mapstructure:"accessKey"`
		SecretKey string `mapstructure:"secretKey"`
		UseSSL    bool   `mapstructure:"useSSL"`
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
	if err := viper.UnmarshalKey("tcp", &conf.Tcp); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("http", &conf.Http); err != nil {
		return err
	}
	if err := envconfig.Process("api", &conf.Api); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("minio", &conf.MinIO); err != nil {
		return err
	}

	return nil
}

func setFromEnv(conf *Config) error {
	if err := envconfig.Process("tcp", &conf.Tcp); err != nil {
		return err
	}
	if err := envconfig.Process("http", &conf.Http); err != nil {
		return err
	}
	if err := envconfig.Process("minio", &conf.MinIO); err != nil {
		return err
	}

	conf.Environment = os.Getenv("APP_ENV")

	return nil
}
