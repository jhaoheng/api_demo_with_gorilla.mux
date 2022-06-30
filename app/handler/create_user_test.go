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

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

/*
[USAGE]
- export AWS_PROFILE=prod
- update the word of `CreateUser`
- run: `go test -run TestCreateUser` or `go test -v ./...`
*/
type SuiteCreateUserTestPlan struct {
	// Account    string
	// Password   string
	// Fullname   string
	ApiMethod  string
	ApiUrl     string
	ApiBody    CreateUserBody
	ExpectCode int
	ExpectBody string
}

type SuiteCreateUser struct {
	suite.Suite
	TestPlans []SuiteCreateUserTestPlan
}

func TestCreateUser(t *testing.T) {
	suite.Run(t, new(SuiteCreateUser))
}

func (s *SuiteCreateUser) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest,", s.T().Name())
	s.TestPlans = []SuiteCreateUserTestPlan{
		0: {
			ApiMethod: "POST",
			ApiUrl:    "/signup",
			ApiBody: CreateUserBody{
				Account:  "max",
				Password: "12345",
				Fullname: "maxhu",
			},
			ExpectCode: http.StatusOK,
			ExpectBody: `{"data":{"account":"max","fullname":"maxhu","created_at":"2022-01-01 12:00:00","updated_at":"2022-01-01 12:00:00"},"error":"0"}`,
		},
		1: {
			ApiMethod: "POST",
			ApiUrl:    "/signup",
			ApiBody: CreateUserBody{
				Account:  "max",
				Password: "*",
				Fullname: "maxhu",
			},
			ExpectCode: http.StatusUnprocessableEntity,
			ExpectBody: `{"data":null,"error":"Key: 'CreateUserBody.Password' Error:Field validation for 'Password' failed on the 'is_allow_password' tag"}`,
		},
		2: {
			ApiMethod: "POST",
			ApiUrl:    "/signup",
			ApiBody: CreateUserBody{
				Account:  "max",
				Password: "123",
				Fullname: "",
			},
			ExpectCode: http.StatusUnprocessableEntity,
			ExpectBody: `{"data":null,"error":"Key: 'CreateUserBody.Fullname' Error:Field validation for 'Fullname' failed on the 'required' tag"}`,
		},
		3: {
			ApiMethod: "POST",
			ApiUrl:    "/signup",
			ApiBody: CreateUserBody{
				Account:  "max",
				Password: "12345",
				Fullname: "maxhu",
			},
			ExpectCode: http.StatusBadRequest,
			ExpectBody: `{"data":null,"error":"Error 1062: Duplicate entry 'account_3' for key 'index'"}`,
		},
		4: {
			ApiMethod:  "POST",
			ApiUrl:     "/signup",
			ApiBody:    CreateUserBody{},
			ExpectCode: http.StatusBadRequest,
			ExpectBody: `{"data":null,"error":"json: cannot unmarshal string into Go value of type handler.CreateUserBody"}`,
		},
	}
	//
	modules.InitValidate()
}

func (s *SuiteCreateUser) TestDo() {
	var err error
	var req *http.Request
	for index, test_plan := range s.TestPlans {
		var api_body_content interface{} = test_plan.ApiBody
		if index == 4 {
			api_body_content = "123"
		}
		req, err = http.NewRequest(test_plan.ApiMethod, test_plan.ApiUrl, func() io.Reader {
			b, _ := json.Marshal(api_body_content)
			return bytes.NewBuffer(b)
		}())
		//
		if !s.NoError(err) {
			s.T().Fatal(err)
		}
		rr := httptest.NewRecorder()

		//
		http.HandlerFunc(NewCreateUser(func() *CreateUser {
			api := CreateUser{
				model_create_user: s.mock_create_user(index, test_plan.ApiBody.Account, test_plan.ApiBody.Password, test_plan.ApiBody.Fullname),
				model_get_user:    s.mock_get_user(index, test_plan.ApiBody.Account, test_plan.ApiBody.Fullname),
			}
			return &api
		}())).ServeHTTP(rr, req)

		//
		// fmt.Println("http status_code=>", rr.Code)
		// fmt.Println("header=>", rr.Header())
		// fmt.Println("body=>", rr.Body.String())
		if rr.Code != test_plan.ExpectCode {
			s.T().Fatalf("handler returned wrong status code: got %v want %v", rr.Code, test_plan.ExpectCode)
		}
		if rr.Body.String() != test_plan.ExpectBody {
			s.T().Fatalf("handler returned unexpected body: \n- got %v \n- want %v", rr.Body.String(), test_plan.ExpectBody)
		}
	}
}

func (s *SuiteCreateUser) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest,", s.T().Name())
}

/*
- mock
*/
func (s *SuiteCreateUser) mock_create_user(index int, acct, pass, fullname string) *models.MockUser {
	mock_create_user := models.NewMockUser()
	mock_create_user.On("SetAcct", acct)
	mock_create_user.On("SetPwd", modules.HashPasswrod(pass))
	mock_create_user.On("SetFullname", fullname)

	if index != 3 {
		mock_create_user.On("Create").Return(nil)
	} else {
		mock_create_user.On("Create").Return(fmt.Errorf("Error 1062: Duplicate entry 'account_3' for key 'index'"))
	}
	return mock_create_user
}

func (s *SuiteCreateUser) mock_get_user(index int, acct, fullname string) *models.MockUser {
	time_at, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 12:00:00")

	mock_get_user := models.NewMockUser()
	mock_get_user.On("SetAcct", acct)
	mock_get_user.On("SetFullname", fullname)
	mock_get_user.On("Get").Return(models.User{
		Acct:      acct,
		Fullname:  fullname,
		CreatedAt: time_at,
		UpdatedAt: time_at,
	}, nil)
	return mock_get_user
}
