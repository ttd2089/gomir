package gomir

import (
	"errors"
	"reflect"
	"testing"
)

func expectErr(t *testing.T, err, expected error) {
	if !errors.Is(err, expected) {
		t.Fatalf("expected <%+v>; actual <%+v>\n", err, expected)
	}
}

func expectExpr(t *testing.T, actual, expected interface{}, expr string) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %s=<%+v>; actual <%+v>\n", expr, expected, actual)
	}
}
