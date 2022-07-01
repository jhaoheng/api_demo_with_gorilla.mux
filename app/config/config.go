package config

import (
	"embed"
	"os"
	"strings"
)

var CFG *ConfigObj

var JWTPubKey []byte
var JWTPriKey []byte

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

//go:embed jwt_rs256.key.pub
//go:embed jwt_rs256.key
var jwt_keypair_embed embed.FS

func LoadJWTKeyPair() {
	keypair, err := NewJWTKeyPair().SetJWTKeypair(jwt_keypair_embed)
	if err != nil {
		panic(err)
	}
	JWTPubKey = keypair.pub_key
	JWTPriKey = keypair.pri_key
}

func NewConfig(env string) *ConfigObj {
	LoadJWTKeyPair()
	//
	CFG = &ConfigObj{}
	switch env {
	default:
		CFG.set_dev()
	}
	return CFG
}

func (c *ConfigObj) set_dev() {
	c.DB_HOST = "localhost"
	if db := os.Getenv("db_host"); len(db) != 0 {
		// docker env
		c.DB_HOST = db
	}
	c.DB_NAME = "my_side_project"
	c.DB_USERNAME = "test"
	c.DB_PASSWORD = "test"
	//
	c.CSRFTOKEN_ONLY_HTTPS = false
	if strings.EqualFold(os.Getenv("csrftoken_only_https"), "true") {
		// docker env
		c.CSRFTOKEN_ONLY_HTTPS = true
	}
	//
	c.JWT_PUBLIC_KEY_PATH = "/tmp/localhost/docker-compose/keypair/jwt_rs256.key.pub"
	c.JWT_PRIVATE_KEY_PATH = "/tmp/localhost/docker-compose/keypair/jwt_rs256.key"
}

// func (c *ConfigObj) set_stg() {
// 	stg := c.get_aws_parameters("/stg")
// 	c = &ConfigObj{
// 		DB_HOST:     stg.DB_HOST,
// 		DB_NAME:     stg.DB_NAME,
// 		DB_USERNAME: stg.DB_USERNAME,
// 		DB_PASSWORD: stg.DB_PASSWORD,
// 		//
// 		CSRFTOKEN_ONLY_HTTPS: true,
// 		//
// 		JWT_PUBLIC_KEY_PATH:  stg.JWT_PUBLIC_KEY_PATH,
// 		JWT_PRIVATE_KEY_PATH: stg.JWT_PRIVATE_KEY_PATH,
// 	}
// }

// func (c *ConfigObj) set_prod() {
// 	prod := c.get_aws_parameters("/prod")
// 	c = &ConfigObj{
// 		DB_HOST:     prod.DB_HOST,
// 		DB_NAME:     prod.DB_NAME,
// 		DB_USERNAME: prod.DB_USERNAME,
// 		DB_PASSWORD: prod.DB_PASSWORD,
// 		//
// 		CSRFTOKEN_ONLY_HTTPS: true,
// 		//
// 		JWT_PUBLIC_KEY_PATH:  prod.JWT_PUBLIC_KEY_PATH,
// 		JWT_PRIVATE_KEY_PATH: prod.JWT_PRIVATE_KEY_PATH,
// 	}
// }

// type AwsParameters struct {
// 	DB_HOST              string
// 	DB_USERNAME          string
// 	DB_PASSWORD          string
// 	DB_NAME              string
// 	JWT_PUBLIC_KEY_PATH  string
// 	JWT_PRIVATE_KEY_PATH string
// }

// func (c *ConfigObj) get_aws_parameters(path string) (aws_parameters *AwsParameters) {
// 	modules.NewAWSSrv()
// 	result, _ := modules.AWSSrv.SSM_GetParametersByPath(path)
// 	for _, parameter := range result.Parameters {
// 		switch *parameter.Name {
// 		case "DB_HOST":
// 			aws_parameters.DB_HOST = *parameter.Value
// 		case "DB_USERNAME":
// 			aws_parameters.DB_USERNAME = *parameter.Value
// 		case "DB_PASSWORD":
// 			aws_parameters.DB_PASSWORD = *parameter.Value
// 		case "DB_NAME":
// 			aws_parameters.DB_NAME = *parameter.Value
// 		case "JWT_PUBLIC_KEY_PATH":
// 			aws_parameters.JWT_PUBLIC_KEY_PATH = *parameter.Value
// 		case "JWT_PRIVATE_KEY_PATH":
// 			aws_parameters.JWT_PRIVATE_KEY_PATH = *parameter.Value
// 		}
// 	}
// 	return
// }
