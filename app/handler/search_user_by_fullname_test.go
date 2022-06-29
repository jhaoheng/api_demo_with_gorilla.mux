package handler

import (
	"app/models"
	"app/modules"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

/*
[USAGE]
- export AWS_PROFILE=prod
- update the word of `SearchUserByFullname`
- run: `go test -run TestSearchUserByFullname` or `go test -v ./...`
*/
type SuiteSearchUserByFullnameTestPlan struct {
	Account    string
	Fullname   string
	ExpectCode int
	ExpectBody string
}

type SuiteSearchUserByFullname struct {
	suite.Suite
	ApiMethod string
	ApiUrl    string
	ApiBody   io.Reader
	TestPlans []SuiteSearchUserByFullnameTestPlan
}

func TestSearchUserByFullname(t *testing.T) {
	suite.Run(t, new(SuiteSearchUserByFullname))
}

func (s *SuiteSearchUserByFullname) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest, ", s.T().Name())
	//
	test_plans := []SuiteSearchUserByFullnameTestPlan{
		0: {
			Account:    "max",
			Fullname:   "maxhu",
			ExpectCode: http.StatusOK,
			ExpectBody: `{"data":{"account":"max","fullname":"maxhu","created_at":"2022-01-01 12:00:00","updated_at":"2022-01-01 12:00:00"},"error":"0"}`,
		},
		1: {
			Account:    "max",
			Fullname:   "maxhu",
			ExpectCode: http.StatusBadRequest,
			ExpectBody: `{"data":null,"error":"mock_get_user error"}`,
		},
	}
	//
	s.ApiMethod = ""
	s.ApiUrl = "/"
	s.ApiBody = nil
	s.TestPlans = test_plans
}

func (s *SuiteSearchUserByFullname) TestDo() {
	for index, test_plan := range s.TestPlans {
		req, err := http.NewRequest(s.ApiMethod, s.ApiUrl, s.ApiBody)
		if !s.NoError(err) {
			s.T().Fatal(err)
		}
		// context.Set(req, "account", test_plan.AccessAccount)
		rr := httptest.NewRecorder()
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			api := SearchUserByFullname{
				path: &SearchUserByFullnamePath{
					Fullname: test_plan.Fullname,
				},
				model_get_user: s.mock_get_user(index, test_plan.Account, test_plan.Fullname),
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
		if rr.Code != test_plan.ExpectCode {
			s.T().Fatalf("handler returned wrong status code: got %v want %v", rr.Code, test_plan.ExpectCode)
		}
		if rr.Body.String() != test_plan.ExpectBody {
			s.T().Fatalf("handler returned unexpected body: \n- got %v \n- want %v", rr.Body.String(), test_plan.ExpectBody)
		}
	}
}

func (s *SuiteSearchUserByFullname) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest, ", s.T().Name())
}

//
func (s *SuiteSearchUserByFullname) mock_get_user(index int, acct, fullname string) *models.MockUser {
	time_at, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 12:00:00")

	mock_get_user := models.NewMockUser()
	mock_get_user.On("SetAcct", acct)
	mock_get_user.On("SetFullname", fullname)

	if index == 0 {
		mock_get_user.On("Get").Return(models.User{
			Acct:      acct,
			Fullname:  fullname,
			CreatedAt: time_at,
			UpdatedAt: time_at,
		}, nil)
	} else {
		mock_get_user.On("Get").Return(models.User{}, fmt.Errorf("mock_get_user error"))
	}
	return mock_get_user
}
