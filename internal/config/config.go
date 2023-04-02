package config

import (
	"fmt"

	"github.com/spf13/viper"
)

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
}

type PrometheusConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Path     string `mapstructure:"path"`
}

type ServicesConfig struct {
	Minio      MinioConfig      `mapstructure:"minio"`
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	ClamAV     ClamAVConfig     `mapstructure:"clamav"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
}

type Config struct {
	CachePath  string         `mapstructure:"cache_path"`
	CachePerms int            `mapstructure:"cache_perms"`
	Services   ServicesConfig `mapstructure:"services"`
}

var vp *viper.Viper

func GetConfig() (*Config, error) {
	vp = viper.New()
	var config Config

	vp.SetConfigName("aegis")
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
