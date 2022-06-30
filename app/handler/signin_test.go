package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api_demo_with_gorilla.mux/app/models"
	"api_demo_with_gorilla.mux/app/modules"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

/*
[USAGE]
- 若有使用 AWS, 則 export AWS_PROFILE=prod
- 更新 `Signin` 為 handler 名稱
- 確定測試資料
- 確定 api => method, path, body
- 更新要替換的 mock
- run: `go test -run TestSignin` or `go test -v ./...`
*/
type SuiteSigninTestPlan struct {
	ApiMethod  string
	ApiUrl     string
	ApiBody    *SigninBody
	ExpectCode int
	ExpectBody string
}

type SuiteSignin struct {
	suite.Suite
	ApiMethod string
	ApiUrl    string
	ApiBody   io.Reader
	TestPlans []SuiteSigninTestPlan
}

func TestSignin(t *testing.T) {
	suite.Run(t, new(SuiteSignin))
}

func (s *SuiteSignin) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest, ", s.T().Name())
	modules.InitValidate()
	//
	test_plans := []SuiteSigninTestPlan{
		0: {
			ApiMethod: "POST",
			ApiUrl:    "/signin",
			ApiBody: &SigninBody{
				Account:  "max",
				Password: "12345",
			},
			ExpectCode: http.StatusOK,
			ExpectBody: "",
		},
		1: {
			ApiMethod:  "POST",
			ApiUrl:     "/signin",
			ApiBody:    &SigninBody{},
			ExpectCode: http.StatusUnprocessableEntity,
			ExpectBody: `{"data":null,"error":"Key: 'SigninBody.Account' Error:Field validation for 'Account' failed on the 'required' tag\nKey: 'SigninBody.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		2: {
			ApiMethod:  "POST",
			ApiUrl:     "/signin",
			ApiBody:    &SigninBody{},
			ExpectCode: http.StatusBadRequest,
			ExpectBody: `{"data":null,"error":"invalid character 'w' looking for beginning of value"}`,
		},
		3: {
			ApiMethod: "POST",
			ApiUrl:    "/signin",
			ApiBody: &SigninBody{
				Account:  "max",
				Password: "12345",
			},
			ExpectCode: http.StatusBadRequest,
			ExpectBody: `{"data":null,"error":"db error"}`,
		},
	}
	s.TestPlans = test_plans
}

func (s *SuiteSignin) TestDo() {
	WSConnections = 1
	for index, test_plan := range s.TestPlans {
		body := func() io.Reader {
			b, _ := json.Marshal(test_plan.ApiBody)
			return bytes.NewBuffer(b)
		}()
		if index == 2 {
			body = bytes.NewBuffer([]byte(`wrong body`))
		}
		req, err := http.NewRequest(test_plan.ApiMethod, test_plan.ApiUrl, body)
		if !s.NoError(err) {
			s.T().Fatal(err)
		}
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/signin", NewSignin(func() *Signin {
			mock_api := Signin{
				model_get_user:       s.mock_get_user(index, test_plan.ApiBody),
				jwt_public_key_path:  "../../keypair/jwt_rs256.key.pub",
				jwt_private_key_path: "../../keypair/jwt_rs256.key",
			}
			return &mock_api
		}()))
		router.ServeHTTP(rr, req)

		//
		// fmt.Println("http status_code=>", rr.Code)
		// fmt.Println("header=>", rr.Header())
		fmt.Println("body=>", rr.Body.String())
		if rr.Code != test_plan.ExpectCode {
			s.T().Fatalf("handler returned wrong status code: got %v want %v", rr.Code, test_plan.ExpectCode)
		}
		if rr.Body.String() != test_plan.ExpectBody {
			if test_plan.ExpectCode == http.StatusOK {
				continue
			}
			s.T().Fatalf("handler returned unexpected body: \n- got %v \n- want %v", rr.Body.String(), test_plan.ExpectBody)
		}
	}
}

func (s *SuiteSignin) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest, ", s.T().Name())
}

//
func (s *SuiteSignin) mock_get_user(index int, body *SigninBody) *models.MockUser {
	time_at, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 12:00:00")

	mock_get_user := models.NewMockUser()
	mock_get_user.On("SetAcct", body.Account)
	mock_get_user.On("SetPwd", modules.HashPasswrod(body.Password))

	switch index {
	case 3:
		mock_get_user.On("Get").Return(models.User{}, fmt.Errorf("db error"))
	default:
		mock_get_user.On("Get").Return(models.User{
			Acct:      body.Account,
			Pwd:       modules.HashPasswrod(body.Password),
			CreatedAt: time_at,
			UpdatedAt: time_at,
		}, nil)
	}
	return mock_get_user
}
