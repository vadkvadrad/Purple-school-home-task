package er

import (
	"fmt"
)

const (
	ErrNotAuthorized        = "user are not authorized"
	ErrWrongUserCredentials = "wrong user credentials"
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %s", msg, err)
}
