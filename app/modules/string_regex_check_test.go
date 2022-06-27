package modules

import "testing"

func Test_CheckRegex(t *testing.T) {
	str := ""
	if err := CheckRegex(str); err != nil {
		t.Fatal(err)
	}
}
