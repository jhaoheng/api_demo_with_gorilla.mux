package modules

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func Test_S3_ListBuckets(t *testing.T) {
	AWSSrv = NewAWSSrv()
	result, err := AWSSrv.S3_ListBuckets()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func Test_SSM_GetParameter(t *testing.T) {
	AWSSrv = NewAWSSrv()
	result, err := AWSSrv.SSM_GetParameter("/name/max")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func Test_SSM_GetParametersByPath(t *testing.T) {
	AWSSrv = NewAWSSrv()
	result, err := AWSSrv.SSM_GetParametersByPath("/name")
	if err != nil {
		t.Fatal(err)
	}
	for _, parameter := range result.Parameters {
		switch *parameter.Name {
		case "/name/max":
			fmt.Println("/name/max =>", *parameter.Value)
		case "/name/max/certificates":
			fmt.Println("/name/max/certificates =>", *parameter.Value)
		}
	}
}
