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
	Account    string
	Password   string
	Fullname   string
	ExpectCode int
	ExpectBody string
}

type SuiteCreateUser struct {
	suite.Suite
	ApiMethod string
	ApiUrl    string
	TestPlans []SuiteCreateUserTestPlan
}

func TestCreateUser(t *testing.T) {
	suite.Run(t, new(SuiteCreateUser))
}

func (s *SuiteCreateUser) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest,", s.T().Name())
	s.ApiMethod = "POST"
	s.ApiUrl = "/signup"
	s.TestPlans = []SuiteCreateUserTestPlan{
		0: {
			Account:    "max",
			Password:   "12345",
			Fullname:   "maxhu",
			ExpectCode: http.StatusOK,
			ExpectBody: `{"data":{"account":"max","fullname":"maxhu","created_at":"2022-01-01 12:00:00","updated_at":"2022-01-01 12:00:00"},"error":"0"}`,
		},
		1: {
			Account:    "max",
			Password:   "*",
			Fullname:   "maxhu",
			ExpectCode: http.StatusUnprocessableEntity,
			ExpectBody: `{"data":null,"error":"Key: 'CreateUserBody.Password' Error:Field validation for 'Password' failed on the 'is_allow_password' tag"}`,
		},
		2: {
			Account:    "max",
			Password:   "123",
			Fullname:   "",
			ExpectCode: http.StatusUnprocessableEntity,
			ExpectBody: `{"data":null,"error":"Key: 'CreateUserBody.Fullname' Error:Field validation for 'Fullname' failed on the 'required' tag"}`,
		},
		3: {
			Account:    "max",
			Password:   "12345",
			Fullname:   "maxhu",
			ExpectCode: http.StatusBadRequest,
			ExpectBody: `{"data":null,"error":"Error 1062: Duplicate entry 'account_3' for key 'index'"}`,
		},
	}
	//
	modules.InitValidate()
}

func (s *SuiteCreateUser) TestDo() {
	for index, test_plan := range s.TestPlans {
		//
		req, err := http.NewRequest(s.ApiMethod, s.ApiUrl, func() io.Reader {
			b, _ := json.Marshal(CreateUserBody{
				Account:  test_plan.Account,
				Password: test_plan.Password,
				Fullname: test_plan.Fullname,
			})
			return bytes.NewBuffer(b)
		}())
		if !s.NoError(err) {
			s.T().Fatal(err)
		}
		rr := httptest.NewRecorder()

		//
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			api := CreateUser{}
			api.model_create_user = s.mock_create_user(index, test_plan.Account, test_plan.Password, test_plan.Fullname)
			api.model_get_user = s.mock_get_user(index, test_plan.Account, test_plan.Fullname)
			payload, status, err := api.do(w, r)
			modules.NewResp(w, r).Set(modules.RespContect{
				Data:   payload,
				Stutus: status,
				Error:  err,
			})
		}).ServeHTTP(rr, req)

		//
		// fmt.Println("http status_code=>", rr.Code)
		// fmt.Println("header=>", rr.Header())
		fmt.Println("body=>", rr.Body.String())
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
