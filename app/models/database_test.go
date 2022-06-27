package models

import "testing"

func TestMain(m *testing.M) {
	NewDBMySQL(DBSet{
		Host:    "localhost",
		User:    "test",
		Pass:    "test",
		DBName:  "my_side_project",
		IsDebug: true,
	})
	m.Run()
}

func Test_NewDBMySQL(t *testing.T) {

}
