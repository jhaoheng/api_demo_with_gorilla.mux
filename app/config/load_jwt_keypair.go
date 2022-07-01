package config

import (
	"embed"
	"fmt"
)

/*
-
*/

type JWTKeyPair struct {
	pub_key []byte
	pri_key []byte
}

func NewJWTKeyPair() *JWTKeyPair {
	return &JWTKeyPair{}
}

func (kp *JWTKeyPair) SetJWTKeypair(keypair embed.FS) (*JWTKeyPair, error) {
	if err := kp.load_from_local(keypair); err == nil {
		return kp, nil
	}
	err := kp.load_from_cloud()
	return kp, err
}

func (kp *JWTKeyPair) load_from_local(keypair embed.FS) error {
	var err error
	if kp.pub_key, err = keypair.ReadFile("keypair/jwt_rs256.key.pub"); err != nil {
		return err
	}
	if kp.pri_key, err = keypair.ReadFile("keypair/jwt_rs256.key"); err != nil {
		return err
	}
	return nil
}

func (kp *JWTKeyPair) load_from_cloud() error {
	return fmt.Errorf("not set keypair")
}
