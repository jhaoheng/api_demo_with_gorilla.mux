package handler

import (
	"app/models"
	"app/modules"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
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
	Expect        string
}

type SuiteDeleteUser struct {
	suite.Suite
	ApiMethod string
	ApiUrl    string
	TestPlan  SuiteDeleteUserTestPlan
}

func TestDeleteUser(t *testing.T) {
	suite.Run(t, new(SuiteDeleteUser))
}

func (s *SuiteDeleteUser) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest")
	//
	test_plan := SuiteDeleteUserTestPlan{
		DelAccount:    "max",
		AccessAccount: "max2",
		Expect:        `{"data":{},"error":"0"}`,
	}
	//
	s.ApiMethod = "DELETE"
	s.ApiUrl = fmt.Sprintf("/account/%v", s.TestPlan.DelAccount)
	s.TestPlan = test_plan
}

func (s *SuiteDeleteUser) TestDo() {
	req, err := http.NewRequest(s.ApiMethod, s.ApiUrl, nil)
	if !s.NoError(err) {
		s.T().Fatal(err)
	}
	context.Set(req, "account", s.TestPlan.AccessAccount)
	rr := httptest.NewRecorder()

	//
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api := DeleteUser{
			path: &DeleteUserPath{
				DelAccount: s.TestPlan.DelAccount,
			},
			access_account: s.TestPlan.AccessAccount,
		}
		api.model_del_account = s.mock_delete_user(s.TestPlan.DelAccount)
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
	if rr.Code != http.StatusOK {
		s.T().Fatalf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	if rr.Body.String() != s.TestPlan.Expect {
		s.T().Fatalf("handler returned unexpected body: \n- got %v \n- want %v", rr.Body.String(), s.TestPlan.Expect)
	}
}

func (s *SuiteDeleteUser) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest")
}

func (s *SuiteDeleteUser) mock_delete_user(acct string) *models.MockUser {
	mock_delete_user := models.NewMockUser()
	mock_delete_user.On("SetAcct", acct)
	mock_delete_user.On("Delete").Return(1, nil)
	return mock_delete_user
}
