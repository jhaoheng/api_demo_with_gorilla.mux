package handler

import (
	"app/models"
	"app/modules"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/context"
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
	Account       string
	Password      string
	Fullname      string
	AccessAccount string
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
	logrus.Info("BeforeTest")
	s.ApiMethod = "POST"
	s.ApiUrl = "/signup"
	s.TestPlans = []SuiteCreateUserTestPlan{
		0: {
			Account:  "max",
			Password: "12345",
			Fullname: "maxhu",
		},
	}
}

func (s *SuiteCreateUser) TestDo() {
	for index, test_plan := range s.TestPlans {
		fmt.Printf("\n=== %v ===\n", index)
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
		context.Set(req, "account", test_plan.AccessAccount)
		rr := httptest.NewRecorder()

		//
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			api := CreateUser{}
			api.model_create_user = s.mock_create_user(test_plan.Account, test_plan.Password, test_plan.Fullname)
			api.model_get_user = s.mock_get_user(test_plan.Account, test_plan.Fullname)
			payload, status, err := api.do(w, r)
			modules.NewResp(w, r).Set(modules.RespContect{
				Data:   payload,
				Stutus: status,
				Error:  err,
			})
		}).ServeHTTP(rr, req)

		//
		fmt.Println("http status_code=>", rr.Code)
		fmt.Println("header=>", rr.Header())
		fmt.Println("body=>", rr.Body.String())
		if rr.Code != http.StatusOK {
			s.T().Error("Not OK")
		}
	}
	fmt.Println("")
}

func (s *SuiteCreateUser) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest")
}

/*
- mock
*/
func (s *SuiteCreateUser) mock_create_user(acct, pass, fullname string) *models.MockUser {
	mock_create_user := models.NewMockUser()
	mock_create_user.On("SetAcct", acct)
	mock_create_user.On("SetPwd", modules.HashPasswrod(pass))
	mock_create_user.On("SetFullname", fullname)
	mock_create_user.On("Create").Return(nil)
	return mock_create_user
}

func (s *SuiteCreateUser) mock_get_user(acct, fullname string) *models.MockUser {
	mock_get_user := models.NewMockUser()
	mock_get_user.On("SetAcct", acct)
	mock_get_user.On("SetFullname", fullname)
	mock_get_user.On("Get").Return(models.User{
		Acct:      acct,
		Fullname:  fullname,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)
	return mock_get_user
}
