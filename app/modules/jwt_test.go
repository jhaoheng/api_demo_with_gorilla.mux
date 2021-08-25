package modules

import (
	"fmt"
	"testing"
)

func Test_JWT(t *testing.T) {

	public_key_path := "../../keypair/jwtRS256.key.pub"
	private_key_path := "../../keypair/jwtRS256.key"

	jwtsrv := NewJWTSrv(public_key_path, private_key_path)
	tokenSting := jwtsrv.Encrtpying("maxhu")

	if account, ok := jwtsrv.Validating(tokenSting); !ok {
		t.Fatal("fail")
	} else {
		fmt.Println(account)
	}
}

func Test_setPublicKey(t *testing.T) {
	path := "../../keypair/jwtRS256.key.pub"
	fmt.Println(setPublicKey(path))
}
