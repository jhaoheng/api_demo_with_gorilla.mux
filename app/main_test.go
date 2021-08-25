package main

import (
	"app/handler"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/context"
	"github.com/sirupsen/logrus"
)

var (
	TestRunAccount       = fmt.Sprintf("test%v%v", time.Now().Minute(), time.Now().Second())
	TestRunPassword      = "test"
	TestRunAuthorization = ""
	TestRunCSRFToken     = ""
)

func TestMain(m *testing.M) {
	logrus.Infof("TestRunAccount: %v, TestRunPassword: %v", TestRunAccount, TestRunPassword)
	m.Run()
}

func Test_route_GetCSRFToken(t *testing.T) {
	testObj := TestObj{
		t:           t,
		Method:      "GET",
		Url:         "/csrf",
		HandlerFunc: http.HandlerFunc(handler.GetCSRFToken),
		Body:        nil,
		//
		ExpectedRespStatusCode: http.StatusOK,
		ExpectedRespBody:       `{"data":"success","error":""}`,
	}
	testObj.DoTestRequest()
}

func Test_route_CreateUser(t *testing.T) {
	b, _ := json.Marshal(handler.CreateUserObj{
		Account:  TestRunAccount,
		Password: TestRunPassword,
		Fullname: "abc",
	})

	testObj := TestObj{
		t:           t,
		Method:      "POST",
		Url:         "/signup",
		HandlerFunc: http.HandlerFunc(handler.CreateUser),
		Body:        bytes.NewBuffer(b),
		//
		ExpectedRespStatusCode: http.StatusOK,
		ExpectedRespBody:       `{"data":"success","error":""}`,
	}
	testObj.DoTestRequest()
}

func Test_route_Signin(t *testing.T) {
	b, _ := json.Marshal(handler.SigninObj{
		Account:  TestRunAccount,
		Password: TestRunPassword,
	})

	testObj := TestObj{
		t:           t,
		Method:      "POST",
		Url:         "/signin",
		HandlerFunc: http.HandlerFunc(handler.Signin),
		Body:        bytes.NewBuffer(b),
		//
		ExpectedRespStatusCode: http.StatusOK,
		ExpectedRespBody:       ``,
	}
	testObj.DoTestRequest()
}

func Test_route_ListAllUsers(t *testing.T) {
	testObj := TestObj{
		t:           t,
		Method:      "GET",
		Url:         "/users",
		HandlerFunc: http.HandlerFunc(handler.ListAllUsers),
		Body:        nil,
		//
		ExpectedRespStatusCode: http.StatusOK,
		ExpectedRespBody:       ``,
	}
	testObj.DoTestRequest()
}

func Test_route_SearchUserByFullname(t *testing.T) {
	testObj := TestObj{
		t:           t,
		Method:      "GET",
		Url:         fmt.Sprintf("/user/fullname/%v", TestRunAccount),
		HandlerFunc: http.HandlerFunc(handler.ListAllUsers),
		Body:        nil,
		//
		ExpectedRespStatusCode: http.StatusOK,
		ExpectedRespBody:       ``,
	}
	testObj.DoTestRequest()
}

func Test_route_GetUserDetailedInfo(t *testing.T) {
	testObj := TestObj{
		t:           t,
		Method:      "GET",
		Url:         "/user/me",
		HandlerFunc: http.HandlerFunc(handler.GetUserDetailedInfo),
		Body:        nil,
		//
		ExpectedRespStatusCode: http.StatusOK,
		ExpectedRespBody:       ``,
	}
	testObj.DoTestRequest()
}

func Test_route_DeleteUser(t *testing.T) {
	testObj := TestObj{
		t:           t,
		Method:      "DELETE",
		Url:         fmt.Sprintf("/user/account/%v", TestRunAccount),
		HandlerFunc: http.HandlerFunc(handler.DeleteUser),
		Body:        nil,
		//
		ExpectedRespStatusCode: http.StatusBadRequest,
		ExpectedRespBody:       ``,
	}
	testObj.DoTestRequest()
}

func Test_route_UpdateUser(t *testing.T) {
	b, _ := json.Marshal(handler.UpdateUserObj{
		Password: "test",
		Fullname: "hello1",
	})

	testObj := TestObj{
		t:           t,
		Method:      "PATCH",
		Url:         "/user/me",
		HandlerFunc: http.HandlerFunc(handler.UpdateUser),
		Body:        bytes.NewBuffer(b),
		//
		ExpectedRespStatusCode: http.StatusOK,
		ExpectedRespBody:       ``,
	}
	testObj.DoTestRequest()
}

func Test_route_UpdateSpecificUserFullname(t *testing.T) {
	b, _ := json.Marshal(handler.UpdateSpecificUserFullnameObj{
		Account:  TestRunAccount,
		Fullname: "hello2",
	})

	testObj := TestObj{
		t:           t,
		Method:      "PATCH",
		Url:         fmt.Sprintf("/user/account/%v", TestRunAccount),
		HandlerFunc: http.HandlerFunc(handler.UpdateSpecificUserFullname),
		Body:        bytes.NewBuffer(b),
		//
		ExpectedRespStatusCode: http.StatusOK,
		ExpectedRespBody:       ``,
	}
	testObj.DoTestRequest()
}

type TestObj struct {
	t *testing.T
	//
	Method      string
	Url         string
	HandlerFunc http.HandlerFunc
	Body        io.Reader
	//
	ExpectedRespStatusCode int
	ExpectedRespBody       string
}

func (testObj *TestObj) DoTestRequest() (header http.Header, respBody string) {
	fmt.Printf("\n【method: %v, url: %v】\n", testObj.Method, testObj.Url)
	//
	t := testObj.t
	//
	req, err := http.NewRequest(testObj.Method, testObj.Url, testObj.Body)
	if err != nil {
		t.Fatal(err)
	}
	context.Set(req, "account", TestRunAccount)
	rr := httptest.NewRecorder()
	testObj.HandlerFunc.ServeHTTP(rr, req)

	if status := rr.Code; status != testObj.ExpectedRespStatusCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	rr.Header()
	return rr.Header(), rr.Body.String()
}
