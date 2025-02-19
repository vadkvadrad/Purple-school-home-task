package er

import (
	"errors"
	"fmt"
)

func Wrap(msg string, err error) error {
	return errors.New(fmt.Sprint(msg, err))
}