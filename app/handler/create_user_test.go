package handler

import (
	"app/models"
	"app/modules"
	"bytes"
	"encoding/json"
	"io"
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
- update the word of `CreateUser`
- run: `go test -run TestCreateUser` or `go test -v ./...`
*/
type SuiteCreateUser struct {
	suite.Suite
}

func TestCreateUser(t *testing.T) {
	suite.Run(t, new(SuiteCreateUser))
}

func (s *SuiteCreateUser) BeforeTest(suiteName, testName string) {
	logrus.Info("BeforeTest")
}

func (s *SuiteCreateUser) TestDo() {
	//
	req, err := http.NewRequest("POST", "/signup", func() io.Reader {
		b, _ := json.Marshal(CreateUserBody{
			Account:  "max",
			Password: "123",
			Fullname: "maxhu",
		})
		return bytes.NewBuffer(b)
	}())
	if !s.NoError(err) {
		s.T().Fatal(err)
	}
	context.Set(req, "account", "max")
	rr := httptest.NewRecorder()

	//
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api := CreateUser{}
		api.model_create_user = s.mock_create_user()
		api.do(w, r)
	}).ServeHTTP(rr, req)

	//
	if status := rr.Code; status != http.StatusOK {
		s.T().Fatalf("Not OK")
	}
	rr.Header()
}

func (s *SuiteCreateUser) AfterTest(suiteName, testName string) {
	logrus.Info("AfterTest")
}

/*
- mock
*/
func (s *SuiteCreateUser) mock_create_user() *models.MockUser {
	mock_create_user := models.NewMockUser()
	mock_create_user.On("SetAcct", "max")
	mock_create_user.On("SetPwd", modules.HashPasswrod("123"))
	mock_create_user.On("SetFullname", "maxhu")
	mock_create_user.On("Create").Return(nil)
	return mock_create_user
}
