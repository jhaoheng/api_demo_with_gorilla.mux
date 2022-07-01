package config

import (
	"fmt"
	"testing"
)

func Test_NewConfig(t *testing.T) {
	c := NewConfig("", nil)
	fmt.Println(c)
}
