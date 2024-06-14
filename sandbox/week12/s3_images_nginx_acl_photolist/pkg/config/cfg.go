package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP struct {
		Port string
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
	vpr := viper.New()
	for key, value := range defaults {
		vpr.SetDefault(key, value)
	}

	vpr.SetConfigName(appName)
	// vpr.AddConfigPath("/etc/")
	vpr.AddConfigPath("week12/s3_images_nginx_acl_photolist/configs/")
	viper.AddConfigPath(".") // optionally look for config in the working directory

	vpr.AutomaticEnv()
	vpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := vpr.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if cfg != nil {
		err := vpr.Unmarshal(cfg)
		if err != nil {
			return nil, err
		}
	}
	return vpr, nil
}
