package models

import (
	"github.com/stretchr/testify/mock"
)

type MockUser struct {
	mock.Mock
	User
}

func NewMockUser() *MockUser {
	return new(MockUser)
}

//
func (mock *MockUser) Or(users ...User) IUser {
	mock.Called(users)
	return mock
}

func (mock *MockUser) SetAcct(acct string) IUser {
	mock.Called(acct)
	return mock
}

func (mock *MockUser) SetPwd(password string) IUser {
	mock.Called(password)
	return mock
}

func (mock *MockUser) SetFullname(fullname string) IUser {
	mock.Called(fullname)
	return mock
}

func (mock *MockUser) Create() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockUser) Get() (User, error) {
	args := mock.Called()
	return args.Get(0).(User), args.Error(1)
}

func (mock *MockUser) GetAll() ([]User, error) {
	args := mock.Called()
	return args.Get(0).([]User), args.Error(1)
}

func (mock *MockUser) GetAllCount() (int64, error) {
	args := mock.Called()
	return int64(args.Int(0)), args.Error(1)
}

func (mock *MockUser) Delete() (int64, error) {
	args := mock.Called()
	return int64(args.Int(0)), args.Error(1)
}

func (mock *MockUser) Update(user User) (int64, error) {
	args := mock.Called(user)
	return int64(args.Int(0)), args.Error(1)
}

func (mock *MockUser) ListBy(paging, sorting string, page_size int) ([]User, error) {
	args := mock.Called(paging, sorting, page_size)
	return args.Get(0).([]User), args.Error(1)
}
