package modules

import (
	"fmt"
	"testing"

	"api_demo_with_gorilla.mux/app/config"
)

func Test_NewJWTSrv(t *testing.T) {
	pub_key := []byte(config.FakePubKey)
	pri_key := []byte(config.FakePriKey)
	jwtsrv, _ := NewJWTSrv(pub_key, pri_key)
	tokenSting := jwtsrv.Encrtpying("maxhu")

	if account, ok := jwtsrv.Validating(tokenSting); !ok {
		t.Fatal("fail")
	} else {
		fmt.Println(account)
	}
}
