package handler

import (
	"bytes"
	"context"
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
- 若有使用 AWS, 則 export AWS_PROFILE=prod
- 更新 `UpdateUser` 為 handler 名稱
- 確定測試資料
- 確定 api => method, path, body
- 更新要替換的 mock
- run: `go test -run TestUpdateUser` or `go test -v ./...`
*/
type SuiteUpdateUserTestPlan struct {
	ApiMethod     string
	ApiUrl        string
	ApiBody       *UpdateUserBody
	AccessAccount string
	ExpectCode    int
	ExpectBody    string
}

type SuiteUpdateUser struct {
	suite.Suite
	ApiMethod string
	ApiUrl    string
	ApiBody   io.Reader
	TestPlans []SuiteUpdateUserTestPlan
}

func TestUpdateUser(t *testing.T) {
	suite.Run(t, new(SuiteUpdateUser))
}

func (s *SuiteUpdateUser) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest, ", s.T().Name())
	modules.InitValidate()
	//
	test_plans := []SuiteUpdateUserTestPlan{
		0: {
			ApiMethod: "PATCH",
			ApiUrl:    "/me",
			ApiBody: &UpdateUserBody{
				Password: "12345",
				Fullname: "maxhu",
			},
			AccessAccount: "max",
			ExpectCode:    http.StatusOK,
			ExpectBody:    `{"data":{"account":"max","fullname":"maxhu","created_at":"2022-01-01 12:00:00","updated_at":"2022-01-01 12:00:00"},"error":"0"}`,
		},
		1: {
			ApiMethod:     "PATCH",
			ApiUrl:        "/me",
			ApiBody:       &UpdateUserBody{},
			AccessAccount: "max",
			ExpectCode:    http.StatusBadRequest,
			ExpectBody:    `{"data":null,"error":"invalid character 'b' looking for beginning of value"}`,
		},
		2: {
			ApiMethod:     "PATCH",
			ApiUrl:        "/me",
			ApiBody:       &UpdateUserBody{},
			AccessAccount: "max",
			ExpectCode:    http.StatusBadRequest,
			ExpectBody:    `{"data":null,"error":"db error"}`,
		},
		3: {
			ApiMethod:     "PATCH",
			ApiUrl:        "/me",
			ApiBody:       &UpdateUserBody{},
			AccessAccount: "max",
			ExpectCode:    http.StatusBadRequest,
			ExpectBody:    `{"data":null,"error":"db error"}`,
		},
	}
	s.TestPlans = test_plans
}

func (s *SuiteUpdateUser) TestDo() {
	for index, test_plan := range s.TestPlans {
		body := func() io.Reader {
			b, _ := json.Marshal(test_plan.ApiBody)
			return bytes.NewBuffer(b)
		}()
		if index == 1 {
			body = bytes.NewBuffer([]byte(`bad body`))
		}
		req, err := http.NewRequest(test_plan.ApiMethod, test_plan.ApiUrl, body)
		if !s.NoError(err) {
			s.T().Fatal(err)
		}
		//
		type AccountType interface{}
		var account_key AccountType = "account"
		var account_value AccountType = test_plan.AccessAccount
		ctx := context.WithValue(req.Context(), account_key, account_value)
		req = req.WithContext(ctx)
		//
		rr := httptest.NewRecorder()
		http.HandlerFunc(NewUpdateUser(func() *UpdateUser {
			mock_api := UpdateUser{
				model_update_user: s.mock_update_user(index, test_plan.AccessAccount, test_plan.ApiBody.Password, test_plan.ApiBody.Fullname),
				model_get_user:    s.mock_get_user(index, test_plan.AccessAccount, test_plan.ApiBody.Fullname),
			}
			return &mock_api
		}())).ServeHTTP(rr, req)

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

func (s *SuiteUpdateUser) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest, ", s.T().Name())
}

//
func (s *SuiteUpdateUser) mock_update_user(index int, acct, password, fullname string) *models.MockUser {
	mock_update_user := models.NewMockUser()
	mock_update_user.On("SetAcct", acct)

	switch index {
	case 2:
		mock_update_user.On("Update", models.User{}).Return(0, fmt.Errorf("db error"))
	default:
		mock_update_user.On("Update", models.User{
			Pwd:      modules.HashPasswrod(password),
			Fullname: fullname,
		}).Return(1, nil)
	}
	return mock_update_user
}

//
func (s *SuiteUpdateUser) mock_get_user(index int, acct, fullname string) *models.MockUser {
	time_at, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 12:00:00")

	mock_get_user := models.NewMockUser()
	mock_get_user.On("SetAcct", acct)

	switch index {
	case 3:
		mock_get_user.On("Get").Return(models.User{}, fmt.Errorf("db error"))
	default:
		mock_get_user.On("Get").Return(models.User{
			Acct:      acct,
			Fullname:  fullname,
			CreatedAt: time_at,
			UpdatedAt: time_at,
		}, nil)
	}
	return mock_get_user
}
