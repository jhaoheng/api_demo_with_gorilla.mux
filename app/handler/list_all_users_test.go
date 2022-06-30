package handler

import (
	"fmt"
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
- 更新 `ListAllUsers` 為 handler 名稱
- 確定測試資料
- 確定 api => method, path, body
- 更新要替換的 mock
- run: `go test -run TestListAllUsers` or `go test -v ./...`
*/
type SuiteListAllUsersTestPlan struct {
	ApiMethod  string
	ApiUrl     string
	ExpectCode int
	ExpectBody string
}

type SuiteListAllUsers struct {
	suite.Suite
	TestPlans []SuiteListAllUsersTestPlan
}

func TestListAllUsers(t *testing.T) {
	suite.Run(t, new(SuiteListAllUsers))
}

func (s *SuiteListAllUsers) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest, ", s.T().Name())
	modules.InitValidate()
	//
	test_plans := []SuiteListAllUsersTestPlan{
		0: {
			ApiMethod:  "GET",
			ApiUrl:     "/users?paging=0",
			ExpectCode: http.StatusOK,
			ExpectBody: `{"data":{"total":1,"users":[{"account":"max","fullname":"maxhu","created_at":"2022-01-01 12:00:00","updated_at":"2022-01-01 12:00:00"}]},"error":"0"}`,
		},
		1: {
			ApiMethod:  "GET",
			ApiUrl:     "/users?paging=hello&sorting=asc",
			ExpectCode: http.StatusBadRequest,
			ExpectBody: `{"data":null,"error":"paging must be number"}`,
		},
		2: {
			ApiMethod:  "GET",
			ApiUrl:     "/users?paging=1&sorting=bbb",
			ExpectCode: http.StatusBadRequest,
			ExpectBody: `{"data":null,"error":"sorting must be 'asc' or 'desc'"}`,
		},
		3: {
			ApiMethod:  "GET",
			ApiUrl:     "/users?paging=0",
			ExpectCode: http.StatusBadRequest,
			ExpectBody: `{"data":null,"error":"test db issue"}`,
		},
		4: {
			ApiMethod:  "GET",
			ApiUrl:     "/users?paging=0",
			ExpectCode: http.StatusBadRequest,
			ExpectBody: `{"data":null,"error":"test db issue"}`,
		},
	}
	s.TestPlans = test_plans
}

func (s *SuiteListAllUsers) TestDo() {
	for index, test_plan := range s.TestPlans {
		req, err := http.NewRequest(test_plan.ApiMethod, test_plan.ApiUrl, nil)
		if !s.NoError(err) {
			s.T().Fatal(err)
		}
		// type AccountType interface{}
		// var account_key AccountType = "account"
		// var account_value AccountType = test_plan.AccessAccount
		// ctx := context.WithValue(req.Context(), account_key, account_value)
		// req = req.WithContext(ctx)
		//
		rr := httptest.NewRecorder()
		//
		router := mux.NewRouter()
		router.HandleFunc("/users", NewListAllUsers(func() *ListAllUsers {
			mock_api := ListAllUsers{
				model_get_all_counts: s.mock_get_all_counts(index),
				model_get_users:      s.mock_get_users(index),
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
			s.T().Fatalf("handler returned unexpected body: \n- got %v \n- want %v", rr.Body.String(), test_plan.ExpectBody)
		}
	}
}

func (s *SuiteListAllUsers) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest, ", s.T().Name())
}

//
func (s *SuiteListAllUsers) mock_get_all_counts(index int) *models.MockUser {
	mock_get_all_count := models.NewMockUser()
	switch index {
	case 3:
		mock_get_all_count.On("GetAllCount").Return(0, fmt.Errorf("test db issue"))
	default:
		mock_get_all_count.On("GetAllCount").Return(1, nil)
	}
	return mock_get_all_count
}

func (s *SuiteListAllUsers) mock_get_users(index int) *models.MockUser {
	time_at, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 12:00:00")
	mock_get_users := models.NewMockUser()

	switch index {
	case 4:
		mock_get_users.On("ListBy", "1", "asc", 10).Return([]models.User{}, fmt.Errorf("test db issue"))
	default:
		mock_get_users.On("ListBy", "1", "asc", 10).Return([]models.User{
			0: {
				Acct:      "max",
				Fullname:  "maxhu",
				CreatedAt: time_at,
				UpdatedAt: time_at,
			},
		}, nil)
	}

	return mock_get_users
}
