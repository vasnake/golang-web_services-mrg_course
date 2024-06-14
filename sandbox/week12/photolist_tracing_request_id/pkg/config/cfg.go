package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP struct {
		Port int
	}
	DB struct {
		Host     string
		Username string
		Password string
		Database string
	}
	S3 struct {
		Host   string
		Access string
		Secret string
		Bucket string
	}
	Session struct {
		Type   string
		Secret string
	}
	Token struct {
		Type   string
		Secret string
	}
	Example struct {
		None string
		Yaml string
		Env1 string
		Env2 string
	}
}

var Defaults = map[string]interface{}{
	"http": map[string]string{
		"port": "8080",
	},
	"db": map[string]string{
		"host":     "dbMysql:3306",
		"username": "root",
		"password": "love",
		"database": "photolist",
	},
	"s3": map[string]string{
		"host":   "http://minio:9000",
		"access": "access_123",
		"secret": "secret_123",
		"bucket": "photolist",
	},
	"session": map[string]string{
		"type":   "jwt_ver",
		"secret": "golangcourseSessionSecret",
	},
	"token": map[string]string{
		"type":   "jwt",
		"secret": "qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB",
	},
}

func Read(appName string, defaults map[string]interface{}, cfg interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(appName)
	v.AddConfigPath("/etc/")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if cfg != nil {
		err := v.Unmarshal(cfg)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}
