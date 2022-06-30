package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"api_demo_with_gorilla.mux/app/models"
	"api_demo_with_gorilla.mux/app/modules"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

/*
[USAGE]
- export AWS_PROFILE=prod
- update the word of `DeleteUser`
- run: `go test -run TestDeleteUser` or `go test -v ./...`
*/
type SuiteDeleteUserTestPlan struct {
	DelAccount    string
	AccessAccount string
	ApiMethod     string
	ApiUrl        string
	ExpectCode    int
	ExpectBody    string
}

type SuiteDeleteUser struct {
	suite.Suite
	TestPlans []SuiteDeleteUserTestPlan
}

func TestDeleteUser(t *testing.T) {
	suite.Run(t, new(SuiteDeleteUser))
}

func (s *SuiteDeleteUser) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest")
	//
	test_plans := []SuiteDeleteUserTestPlan{
		0: {
			DelAccount:    "max",
			AccessAccount: "max2",
			ApiMethod:     "DELETE",
			ApiUrl:        fmt.Sprintf("/account/%v", "max"),
			ExpectCode:    http.StatusOK,
			ExpectBody:    `{"data":{},"error":"0"}`,
		},
		1: {
			DelAccount:    "max",
			AccessAccount: "max2",
			ApiMethod:     "DELETE",
			ApiUrl:        fmt.Sprintf("/account/%v", "max"),
			ExpectCode:    http.StatusBadRequest,
			ExpectBody:    `{"data":null,"error":"test db error"}`,
		},
	}
	s.TestPlans = test_plans
	//
	modules.InitValidate()
}

func (s *SuiteDeleteUser) TestDo() {
	for index, test_plan := range s.TestPlans {
		req, err := http.NewRequest(test_plan.ApiMethod, test_plan.ApiUrl, nil)
		if !s.NoError(err) {
			s.T().Fatal(err)
		}
		type AccountType interface{}
		var account_key AccountType = "account"
		var account_value AccountType = test_plan.AccessAccount
		ctx := context.WithValue(req.Context(), account_key, account_value)
		req = req.WithContext(ctx)
		//
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/account/{account}", NewDeleteUser(func() *DeleteUser {
			api := DeleteUser{
				model_del_account: s.mock_delete_user(index, test_plan.DelAccount),
			}
			return &api
		}()))
		router.ServeHTTP(rr, req)

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

func (s *SuiteDeleteUser) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest")
}

func (s *SuiteDeleteUser) mock_delete_user(index int, acct string) *models.MockUser {
	mock_delete_user := models.NewMockUser()
	mock_delete_user.On("SetAcct", acct)
	switch index {
	case 1:
		mock_delete_user.On("Delete").Return(1, fmt.Errorf("test db error"))
	default:
		mock_delete_user.On("Delete").Return(1, nil)
	}
	return mock_delete_user
}
