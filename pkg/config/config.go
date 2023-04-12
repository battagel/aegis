package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	LoggerLevel    string `mapstructure:"LOGGER_LEVEL"`
	LoggerEncoding string `mapstructure:"LOGGER_ENCODING"`

	MinioEndpoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinioUseSSL    bool   `mapstructure:"MINIO_USE_SSL"`

	KafkaBrokers  []string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopic    string   `mapstructure:"KAFKA_TOPIC"`
	KafkaGroupID  string   `mapstructure:"KAFKA_GROUP_ID"`
	KafkaMaxBytes int      `mapstructure:"KAFKA_MAX_BYTES"`

	ClamAVRemoveAfterScan bool   `mapstructure:"CLAMAV_REMOVE_AFTER_SCAN"`
	ClamAVDateTimeFormat  string `mapstructure:"CLAMAV_DATETIME_FORMAT"`
	ClamAVPath            string `mapstructure:"CLAMAV_PATH"`

	PrometheusEndpoint string `mapstructure:"PROMETHEUS_ENDPOINT"`
	PrometheusPath     string `mapstructure:"PROMETHEUS_PATH"`

	PostgresUsername string `mapstructure:"POSTGRES_USERNAME"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresEndpoint string `mapstructure:"POSTGRES_ENDPOINT"`
	PostgresDatabase string `mapstructure:"POSTGRES_DATABASE"`
	PostgresTable    string `mapstructure:"POSTGRES_TABLE"`
}

var vp *viper.Viper

func GetConfig() (*Config, error) {
	vp = viper.New()
	var config Config

	vp.SetConfigName("config")
	vp.SetConfigType("env")
	vp.AddConfigPath(".")

	err := vp.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return &Config{}, err
	}

	// Read in environment variables
	vp.AutomaticEnv()

	err = vp.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct: ", err)
		return &Config{}, err
	}
	return &config, nil
}
