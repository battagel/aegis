package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

type KafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	Topic    string   `mapstructure:"topic"`
	GroupID  string   `mapstructure:"group_id"`
	MaxBytes int      `mapstructure:"max_bytes"`
}

type ClamAVConfig struct {
	RemoveAfterScan bool   `mapstructure:"remove_after_scan"`
	DatetimeFormat  string `mapstructure:"datetime_format"`
	Path            string `mapstructure:"path"`
	Perms           int    `mapstructure:"perms"`
}

type PrometheusConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Path     string `mapstructure:"path"`
}

type PostgresConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Endpoint string `mapstructure:"endpoint"`
	Database string `mapstructure:"database"`
	Table    string `mapstructure:"table"`
}

type Config struct {
	Logger     LoggerConfig     `mapstructure:"logger"`
	Minio      MinioConfig      `mapstructure:"minio"`
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	ClamAV     ClamAVConfig     `mapstructure:"clamav"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	Postgres   PostgresConfig   `mapstructure:"postgres"`
}

var vp *viper.Viper

func GetConfig() (*Config, error) {
	vp = viper.New()
	var config Config

	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")

	err := vp.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return &Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct: ", err)
		return &Config{}, err
	}

	return &config, nil
}
