package handler

import (
	"bytes"
	"encoding/json"
	"io"
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
- 更新 `UpdateUserFullname` 為 handler 名稱
- 確定測試資料
- 確定 api => method, path, body
- 更新要替換的 mock
- run: `go test -run TestUpdateUserFullname` or `go test -v ./...`
*/
type SuiteUpdateUserFullnameTestPlan struct {
	ApiMethod   string
	ApiUrl      string
	ApiBody     *UpdateUserFullnameBody
	TestAccount string
	ExpectCode  int
	ExpectBody  string
}

type SuiteUpdateUserFullname struct {
	suite.Suite
	ApiMethod string
	ApiUrl    string
	ApiBody   io.Reader
	TestPlans []SuiteUpdateUserFullnameTestPlan
}

func TestUpdateUserFullname(t *testing.T) {
	suite.Run(t, new(SuiteUpdateUserFullname))
}

func (s *SuiteUpdateUserFullname) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest, ", s.T().Name())
	modules.InitValidate()
	//
	test_plans := []SuiteUpdateUserFullnameTestPlan{
		0: {
			ApiMethod: "PATCH",
			ApiUrl:    "/account/max",
			ApiBody: &UpdateUserFullnameBody{
				Fullname: "maxhu",
			},
			TestAccount: "max",
			ExpectCode:  http.StatusOK,
			ExpectBody:  `{"data":{"account":"max","fullname":"","created_at":"2022-01-01 12:00:00","updated_at":"2022-01-01 12:00:00"},"error":"0"}`,
		},
	}
	s.TestPlans = test_plans
}

func (s *SuiteUpdateUserFullname) TestDo() {
	for index, test_plan := range s.TestPlans {
		req, err := http.NewRequest(test_plan.ApiMethod, test_plan.ApiUrl, func() io.Reader {
			b, _ := json.Marshal(test_plan.ApiBody)
			return bytes.NewBuffer(b)
		}())
		if !s.NoError(err) {
			s.T().Fatal(err)
		}
		// context.Set(req, "account", test_plan.AccessAccount)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/account/{account}", NewUpdateUserFullname(func() *UpdateUserFullname {
			mock_api := UpdateUserFullname{
				model_update_user: s.mock_update_user(index, test_plan.TestAccount, test_plan.ApiBody.Fullname),
				model_get_user:    s.mock_get_user(index, test_plan.TestAccount),
			}
			return &mock_api
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

func (s *SuiteUpdateUserFullname) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest, ", s.T().Name())
}

//
func (s *SuiteUpdateUserFullname) mock_update_user(index int, acct, fullname string) *models.MockUser {
	mock_update_user := models.NewMockUser()
	mock_update_user.On("SetAcct", acct)
	mock_update_user.On("Update", models.User{
		Fullname: fullname,
	}).Return(1, nil)
	return mock_update_user
}

//
func (s *SuiteUpdateUserFullname) mock_get_user(index int, acct string) *models.MockUser {
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
