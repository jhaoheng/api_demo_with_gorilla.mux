package config

import (
	"app/modules"
	"os"
	"strings"
)

var CFG *ConfigObj

type ConfigObj struct {
	DB_HOST     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME     string
	//
	CSRFTOKEN_ONLY_HTTPS bool
	//
	JWT_PUBLIC_KEY_PATH  string
	JWT_PRIVATE_KEY_PATH string
}

func NewConfig(env string) *ConfigObj {
	CFG = &ConfigObj{}
	switch env {
	case "stg":
		CFG.setStg()
	case "prod":
		CFG.setProd()
	default:
		CFG.setDev()
	}
	return CFG
}

func (c *ConfigObj) setDev() {
	c.DB_HOST = "localhost"
	if db := os.Getenv("db_host"); len(db) != 0 {
		// docker env
		c.DB_HOST = db
	}
	c.DB_NAME = "ui_test"
	c.DB_USERNAME = "ui_test"
	c.DB_PASSWORD = "ui_test"
	//
	c.CSRFTOKEN_ONLY_HTTPS = false
	if strings.EqualFold(os.Getenv("csrftoken_only_https"), "true") {
		// docker env
		c.CSRFTOKEN_ONLY_HTTPS = true
	}
	//
	c.JWT_PUBLIC_KEY_PATH = "/tmp/localhost/docker-compose/keypair/jwtRS256.key.pub"
	c.JWT_PRIVATE_KEY_PATH = "/tmp/localhost/docker-compose/keypair/jwtRS256.key"
}

func (c *ConfigObj) setStg() {
	stg := c.getAwsParameters("/stg")
	c = &ConfigObj{
		DB_HOST:     stg.DB_HOST,
		DB_NAME:     stg.DB_NAME,
		DB_USERNAME: stg.DB_USERNAME,
		DB_PASSWORD: stg.DB_PASSWORD,
		//
		CSRFTOKEN_ONLY_HTTPS: true,
		//
		JWT_PUBLIC_KEY_PATH:  stg.JWT_PUBLIC_KEY_PATH,
		JWT_PRIVATE_KEY_PATH: stg.JWT_PRIVATE_KEY_PATH,
	}
}

func (c *ConfigObj) setProd() {
	prod := c.getAwsParameters("/prod")
	c = &ConfigObj{
		DB_HOST:     prod.DB_HOST,
		DB_NAME:     prod.DB_NAME,
		DB_USERNAME: prod.DB_USERNAME,
		DB_PASSWORD: prod.DB_PASSWORD,
		//
		CSRFTOKEN_ONLY_HTTPS: true,
		//
		JWT_PUBLIC_KEY_PATH:  prod.JWT_PUBLIC_KEY_PATH,
		JWT_PRIVATE_KEY_PATH: prod.JWT_PRIVATE_KEY_PATH,
	}
}

type AwsParameters struct {
	DB_HOST              string
	DB_USERNAME          string
	DB_PASSWORD          string
	DB_NAME              string
	JWT_PUBLIC_KEY_PATH  string
	JWT_PRIVATE_KEY_PATH string
}

func (c *ConfigObj) getAwsParameters(path string) (awsParameters *AwsParameters) {
	modules.NewAWSSrv()
	result, _ := modules.AWSSrv.SSM_GetParametersByPath(path)
	for _, parameter := range result.Parameters {
		switch *parameter.Name {
		case "DB_HOST":
			awsParameters.DB_HOST = *parameter.Value
		case "DB_USERNAME":
			awsParameters.DB_USERNAME = *parameter.Value
		case "DB_PASSWORD":
			awsParameters.DB_PASSWORD = *parameter.Value
		case "DB_NAME":
			awsParameters.DB_NAME = *parameter.Value
		case "JWT_PUBLIC_KEY_PATH":
			awsParameters.JWT_PUBLIC_KEY_PATH = *parameter.Value
		case "JWT_PRIVATE_KEY_PATH":
			awsParameters.JWT_PRIVATE_KEY_PATH = *parameter.Value
		}
	}
	return
}
