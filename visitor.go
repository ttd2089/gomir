package gomir

import "reflect"

// A Visitor recursively iterates the fields of arbitrary values.
type Visitor struct {
}

// Describe describes an arbitrary value's fields using a default Visitor.
func Describe(target interface{}) ([]Property, error) {
	return Visitor{}.Describe(target)
}

// Describe describes an arbitrary value's fields.
func (v Visitor) Describe(target interface{}) ([]Property, error) {
	typ, val, err := indirect(target)
	if err != nil {
		return nil, err
	}
	vv := visitor{
		Visitor: v,
	}
	return vv.describe(typ, val)
}

type visitor struct {
	Visitor
}

func (v visitor) describe(typ reflect.Type, val reflect.Value) ([]Property, error) {
	nFields := typ.NumField()
	props := make([]Property, 0, nFields)
	for i := 0; i < nFields; i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		prop := Property{
			Field: field,
			Value: valueField(val, i),
		}
		fTyp, fVal := indirectProperty(prop)
		if fTyp.Kind() == reflect.Struct {
			props2, err := v.describe(fTyp, fVal)
			if err != nil {
				return nil, err
			}
			prop.Properties = props2
		}
		props = append(props, prop)
	}
	return props, nil
}

func valueField(val reflect.Value, i int) reflect.Value {
	if val.IsValid() {
		return val.Field(i)
	}
	return reflect.Value{}
}

func indirectProperty(prop Property) (reflect.Type, reflect.Value) {
	typ, val := prop.Field.Type, prop.Value
	if typ.Kind() == reflect.Pointer {
		return typ.Elem(), val.Elem()
	}
	return typ, val
}

func indirect(target interface{}) (reflect.Type, reflect.Value, error) {
	p := reflect.ValueOf(target)
	if p.Kind() == reflect.Pointer {
		if val := p.Elem(); val.Kind() == reflect.Struct {
			return val.Type(), val, nil
		}
	}
	return nil, reflect.Value{}, newErrInvalidTarget("must be a non-nil pointer to a struct")
}
