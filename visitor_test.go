package gomir

import (
	"fmt"
	"reflect"
	"testing"
)

func TestVisitor_Describe(t *testing.T) {

	for _, tt := range []struct {
		objType string
		obj     interface{}
	}{
		{"a struct value", struct{}{}},
		{"a scalar value", 13},
		{"a slice", []int{1, 3}},
		{"a map", map[string]string{}},
		{"a chan", make(chan int, 13)},
		{"an array", [13]int{}},
		{"a pointer to a non-struct", new(int)},
	} {
		t.Run(fmt.Sprintf("returns an error when obj is %s", tt.objType), func(t *testing.T) {
			v := Visitor{}
			_, err := v.Describe(tt.obj)
			expectErr(t, err, ErrInvalidTarget)
		})
	}

	for _, tt := range []struct {
		fieldType string
		target    interface{}
	}{
		{
			fieldType: "scalar",
			target:    &struct{ X int }{7},
		},
		{
			fieldType: "slice",
			target:    &struct{ X []int }{[]int{1, 2, 3}},
		},
		{
			fieldType: "map",
			target:    &struct{ X map[int]string }{map[int]string{1: "one"}},
		},
		{
			fieldType: "pointer",
			target:    &struct{ X *int }{new(int)},
		},
	} {
		t.Run(fmt.Sprintf("returns directly owned %s fields", tt.fieldType), func(t *testing.T) {
			v := Visitor{}
			expected := Property{
				Field: reflect.TypeOf(tt.target).Elem().Field(0),
				Value: reflect.ValueOf(tt.target).Elem().Field(0),
			}
			props, err := v.Describe(tt.target)
			expectErr(t, err, nil)
			expectExpr(t, props[0], expected, "props[0]")
		})
	}

	t.Run("returns nested properties of a struct", func(t *testing.T) {
		type Foo struct {
			X int
		}
		type Bar struct {
			Foo Foo
		}
		v := Visitor{}
		target := &Bar{}
		expected := Property{
			Field: reflect.TypeOf(target).Elem().Field(0),
			Value: reflect.ValueOf(target).Elem().Field(0),
			Properties: []Property{
				{
					Field: reflect.TypeOf(target).Elem().FieldByIndex([]int{0, 0}),
					Value: reflect.ValueOf(target).Elem().FieldByIndex([]int{0, 0}),
				},
			},
		}
		props, err := v.Describe(target)
		expectErr(t, err, nil)
		expectExpr(t, props[0], expected, "props[0]")
	})

	t.Run("returns nested properties of a pointer to a struct", func(t *testing.T) {
		type Foo struct {
			X int
		}
		type Bar struct {
			Foo *Foo
		}
		v := Visitor{}
		target := &Bar{}
		expected := Property{
			Field: reflect.TypeOf(target).Elem().Field(0),
			Value: reflect.ValueOf(target).Elem().Field(0),
			Properties: []Property{
				{
					Field: reflect.TypeOf(target).Elem().FieldByIndex([]int{0, 0}),
					Value: reflect.Value{},
				},
			},
		}
		props, err := v.Describe(target)
		expectErr(t, err, nil)
		expectExpr(t, props[0], expected, "props[0]")
	})

	t.Run("returns dereferenced values of nested properties of a pointer to a struct", func(t *testing.T) {
		type Foo struct {
			X int
		}
		type Bar struct {
			Foo *Foo
		}
		type Baz struct {
			Bar *Bar
		}
		v := Visitor{}
		target := &Baz{&Bar{&Foo{7}}}
		expected := Property{
			Field: reflect.TypeOf(target).Elem().Field(0),
			Value: reflect.ValueOf(target).Elem().Field(0),
			Properties: []Property{
				{
					Field: reflect.TypeOf(target).Elem().FieldByIndex([]int{0, 0}),
					Value: reflect.ValueOf(target).Elem().FieldByIndex([]int{0, 0}),
					Properties: []Property{
						{
							Field: reflect.TypeOf(target).Elem().FieldByIndex([]int{0, 0, 0}),
							Value: reflect.ValueOf(target).Elem().FieldByIndex([]int{0, 0, 0}),
						},
					},
				},
			},
		}
		props, err := v.Describe(target)
		expectErr(t, err, nil)
		expectExpr(t, props[0], expected, "props[0]")
	})

	t.Run("skips unexported properties", func(t *testing.T) {
		type Foo struct {
			X int
			y int
		}
		type Bar struct {
			z   int
			Foo Foo
		}
		v := Visitor{}
		target := &Bar{}
		expected := Property{
			Field: reflect.TypeOf(target).Elem().Field(1),
			Value: reflect.ValueOf(target).Elem().Field(1),
			Properties: []Property{
				{
					Field: reflect.TypeOf(target).Elem().FieldByIndex([]int{1, 0}),
					Value: reflect.ValueOf(target).Elem().FieldByIndex([]int{1, 0}),
				},
			},
		}
		props, err := v.Describe(target)
		expectErr(t, err, nil)
		expectExpr(t, len(props), 1, "len(props)")
		expectExpr(t, props[0], expected, "props[0]")
	})
}
