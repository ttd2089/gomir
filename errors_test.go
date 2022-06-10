package gomir

import "testing"

func Test_newErrInvalidTarget(t *testing.T) {
	expectErr(t, newErrInvalidTarget("msg"), ErrInvalidTarget)
}
