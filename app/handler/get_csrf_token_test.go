package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api_demo_with_gorilla.mux/app/modules"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

/*
[USAGE]
- 若有使用 AWS, 則 export AWS_PROFILE=prod
- 更新 `GetCSRFToken` 為 handler 名稱
- 確定測試資料
- 確定 api => method, path, body
- 更新要替換的 mock
- run: `go test -run TestGetCSRFToken` or `go test -v ./...`
*/
type SuiteGetCSRFTokenTestPlan struct {
	ApiMethod  string
	ApiUrl     string
	ExpectCode int
	ExpectBody string
}

type SuiteGetCSRFToken struct {
	suite.Suite
	TestPlans []SuiteGetCSRFTokenTestPlan
}

func TestGetCSRFToken(t *testing.T) {
	suite.Run(t, new(SuiteGetCSRFToken))
}

func (s *SuiteGetCSRFToken) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest, ", s.T().Name())
	modules.InitValidate()
	//
	test_plans := []SuiteGetCSRFTokenTestPlan{
		0: {
			ApiMethod:  "GET",
			ApiUrl:     "/",
			ExpectCode: http.StatusOK,
			ExpectBody: `{"data":null,"error":"0"}`,
		},
	}
	s.TestPlans = test_plans
}

func (s *SuiteGetCSRFToken) TestDo() {
	for _, test_plan := range s.TestPlans {
		req, err := http.NewRequest(test_plan.ApiMethod, test_plan.ApiUrl, nil)
		if !s.NoError(err) {
			s.T().Fatal(err)
		}
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/", NewGetCSRFToken(nil))
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

func (s *SuiteGetCSRFToken) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest, ", s.T().Name())
}
