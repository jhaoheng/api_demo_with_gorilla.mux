package modules

import (
	"errors"
	"regexp"
)

func CheckRegex(str ...string) error {
	var re = regexp.MustCompile(`^[A-Za-z0-9]{1,10}$`)
	for _, s := range str {
		if !re.MatchString(s) {
			return errors.New(`all parameters MUST match '^[A-Za-z0-9]{1,10}$'`)
		}
	}
	return nil
}
