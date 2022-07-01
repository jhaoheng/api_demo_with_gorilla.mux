package modules

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type JWTSRV struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

var JWTSrvAgent *JWTSRV

func NewJWTSrv(pub_key, pri_key []byte) (*JWTSRV, error) {
	if JWTSrvAgent != nil {
		return JWTSrvAgent, nil
	}

	return &JWTSRV{
		publicKey:  set_public_key(pub_key),
		privateKey: set_private_key(pri_key),
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

func set_public_key(key []byte) (public_key *rsa.PublicKey) {
	var err error
	public_key, err = jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		logrus.Fatal(err)
	}
	return
}

func set_private_key(key []byte) (private_key *rsa.PrivateKey) {
	var err error
	private_key, err = jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		logrus.Fatal(err)
	}
	return
}
