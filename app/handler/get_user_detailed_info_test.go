package handler

import (
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
- update the word of `GetUserDetailed`
- run: `go test -run TestGetUserDetailed` or `go test -v ./...`
*/
type SuiteGetUserDetailedTestPlan struct {
	Account    string
	ExpectCode int
	ExpectBody string
}

type SuiteGetUserDetailed struct {
	suite.Suite
	ApiMethod string
	ApiUrl    string
	TestPlan  SuiteGetUserDetailedTestPlan
}

func TestGetUserDetailed(t *testing.T) {
	suite.Run(t, new(SuiteGetUserDetailed))
}

func (s *SuiteGetUserDetailed) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest, ", s.T().Name())
	test_plan := SuiteGetUserDetailedTestPlan{
		Account:    "max",
		ExpectCode: http.StatusOK,
		ExpectBody: `{"data":{"account":"max","fullname":"","created_at":"2022-01-01 12:00:00","updated_at":"2022-01-01 12:00:00"},"error":"0"}`,
	}
	s.ApiMethod = "GET"
	s.ApiUrl = "/"
	s.TestPlan = test_plan
}

func (s *SuiteGetUserDetailed) TestDo() {
	req, err := http.NewRequest(s.ApiMethod, s.ApiUrl, nil)
	if !s.NoError(err) {
		s.T().Fatal(err)
	}
	rr := httptest.NewRecorder()

	//
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api := GetUserDetailed{
			access_account: s.TestPlan.Account,
			model_get_user: s.mock_get_user(s.TestPlan.Account),
		}
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
	// fmt.Println("body=>", rr.Body.String())
	if rr.Code != s.TestPlan.ExpectCode {
		s.T().Fatalf("handler returned wrong status code: got %v want %v", rr.Code, s.TestPlan.ExpectCode)
	}
	if rr.Body.String() != s.TestPlan.ExpectBody {
		s.T().Fatalf("handler returned unexpected body: \n- got %v \n- want %v", rr.Body.String(), s.TestPlan.ExpectBody)
	}
}

func (s *SuiteGetUserDetailed) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest, ", s.T().Name())
}

func (s *SuiteGetUserDetailed) mock_get_user(acct string) *models.MockUser {
	time_at, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 12:00:00")

	mock_get_user := models.NewMockUser()
	mock_get_user.On("SetAcct", acct)
	mock_get_user.On("Get").Return(models.User{
		Acct:      acct,
		Fullname:  "",
		CreatedAt: time_at,
		UpdatedAt: time_at,
	}, nil)
	return mock_get_user
}
