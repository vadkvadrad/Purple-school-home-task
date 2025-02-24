package er

import (
	"fmt"
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %s", msg, err)
}