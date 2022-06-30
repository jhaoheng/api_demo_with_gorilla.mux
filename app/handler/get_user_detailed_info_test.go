package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api_demo_with_gorilla.mux/app/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

/*
[USAGE]
- export AWS_PROFILE=prod
- update the word of `GetUserDetailed`
- run: `go test -run TestGetUserDetailed` or `go test -v ./...`
*/
type SuiteGetUserDetailedTestPlan struct {
	AccessAccount string
	ApiMethod     string
	ApiUrl        string
	ExpectCode    int
	ExpectBody    string
}

type SuiteGetUserDetailed struct {
	suite.Suite
	TestPlans []SuiteGetUserDetailedTestPlan
}

func TestGetUserDetailed(t *testing.T) {
	suite.Run(t, new(SuiteGetUserDetailed))
}

func (s *SuiteGetUserDetailed) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest, ", s.T().Name())
	test_plans := []SuiteGetUserDetailedTestPlan{
		0: {
			AccessAccount: "max",
			ApiMethod:     "GET",
			ApiUrl:        "/",
			ExpectCode:    http.StatusOK,
			ExpectBody:    `{"data":{"account":"max","fullname":"","created_at":"2022-01-01 12:00:00","updated_at":"2022-01-01 12:00:00"},"error":"0"}`,
		},
		1: {
			AccessAccount: "max",
			ApiMethod:     "GET",
			ApiUrl:        "/",
			ExpectCode:    http.StatusBadRequest,
			ExpectBody:    `{"data":null,"error":"test db fail"}`,
		},
	}
	s.TestPlans = test_plans
}

func (s *SuiteGetUserDetailed) TestDo() {
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

		//
		router := mux.NewRouter()
		router.HandleFunc("/", NewGetUserDetailed(func() *GetUserDetailed {
			api := GetUserDetailed{
				model_get_user: s.mock_get_user(index, test_plan.AccessAccount),
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

func (s *SuiteGetUserDetailed) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest, ", s.T().Name())
}

func (s *SuiteGetUserDetailed) mock_get_user(index int, acct string) *models.MockUser {
	time_at, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 12:00:00")

	mock_get_user := models.NewMockUser()
	mock_get_user.On("SetAcct", acct)

	switch index {
	case 1:
		mock_get_user.On("Get").Return(models.User{}, fmt.Errorf("test db fail"))
	default:
		mock_get_user.On("Get").Return(models.User{
			Acct:      acct,
			Fullname:  "",
			CreatedAt: time_at,
			UpdatedAt: time_at,
		}, nil)
	}
	return mock_get_user
}
