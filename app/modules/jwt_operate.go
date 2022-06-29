package modules

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// var (
// 	default_public_key_path  = "../../keypair/jwt_rs256.key.pub"
// 	default_private_key_path = "../../keypair/jwt_rs256.key"
// )

type JWTSRV struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

var JWTSrv *JWTSRV

func NewJWTSrv(pubKeyPath, priKeyPath string) (*JWTSRV, error) {
	if JWTSrv != nil {
		return JWTSrv, nil
	}
	if len(pubKeyPath) == 0 {
		return nil, fmt.Errorf("public_key_path is empty")
	}
	if len(priKeyPath) == 0 {
		return nil, fmt.Errorf("private_key_path is empty")
	}

	return &JWTSRV{
		publicKey:  setPublicKey(pubKeyPath),
		privateKey: setPrivateKey(priKeyPath),
	}, nil
}

// Encrtpying ...
func (j *JWTSRV) Encrtpying(account string) (tokenString string) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Issuer:    "test",
		Subject:   account,
	})

	var err error
	tokenString, err = token.SignedString(j.privateKey)
	if err != nil {
		logrus.Fatal(err)
	}
	return
}

// Validating ...
func (j *JWTSRV) Validating(tokenString string) (account string, ok bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})
	if err != nil {
		logrus.Error(err)
		return "", false
	}

	if !token.Valid {
		logrus.Info("jwt token invalid")
		return "", false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logrus.Info("jwt claims issue")
		return "", false
	}

	return claims["sub"].(string), true
}

func getKeyFromPath(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatal(err)
	}
	return bytes
}

func setPublicKey(public_key_path string) (public_key *rsa.PublicKey) {
	var err error
	public_key, err = jwt.ParseRSAPublicKeyFromPEM(getKeyFromPath(public_key_path))
	if err != nil {
		logrus.Fatal(err)
	}
	return
}

func setPrivateKey(private_key_path string) (private_key *rsa.PrivateKey) {
	var err error
	private_key, err = jwt.ParseRSAPrivateKeyFromPEM(getKeyFromPath(private_key_path))
	if err != nil {
		logrus.Fatal(err)
	}
	return
}