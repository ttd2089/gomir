package gomir

import (
	"errors"
	"fmt"
)

// ErrInvalidTarget is returned when gomir.Describe() or (Visitor).Describe() is called with
// invalid value for target. By default any pointer to struct is valid and anything else is
// invalid. See the options on the Visitor type to configure custom validity checks.
var ErrInvalidTarget = errors.New("invalid target")

func newErrInvalidTarget(msg string) error {
	return fmt.Errorf("%w: %s", ErrInvalidTarget, msg)
}
