package gomir

import "reflect"

// A Property represents the definition and value of a field belonging to an instance of a struct.
// A Property represent a field belonging directly to the struct, embedded within the struct, or
// belonging to a direct child or other descendant of the struct.
type Property struct {

	// Field is the metadata describing the Property.
	Field reflect.StructField

	// Value is the current value of the Property. Value.IsNil() will be true when the Property is
	// holding nil and belongs to a non-nil value within the struct. Value == nil will be true when
	// the Property is a descendent of a nil value within the struct.
	Value reflect.Value

	// Properties are the direct children of the Property.
	Properties []Property
}
