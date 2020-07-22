package gdi

import (
	"fmt"
	"reflect"
)

type resolver interface {
	resolve(b *Binding) error
}

type resolverFunc func(parent, child *Binding, fieldName string)

func (fn resolverFunc) resolve(parent, child *Binding, fieldName string) { fn(parent, child, fieldName) }

func transient(parent, child *Binding, fieldName string) error {
	field := parent.reflectVal.Elem().FieldByName(fieldName)

	if !field.CanSet() {
		return fmt.Errorf(
			"gdi: inject on unexported field %s in object of type %s",
			fieldName,
			parent.reflectTyp,
		)
	}

	if field.Kind() == reflect.Struct && child.reflectTyp.Kind() == reflect.Ptr {
		child.reflectVal = child.reflectVal.Elem()
		child.reflectTyp = child.reflectVal.Type()
	}

	if !field.Type().AssignableTo(child.reflectTyp) {
		return fmt.Errorf(
			"XXXX",
			// "gdi: inject into %s of type %s is not assignable to field %s of type %s in object of type %s",
			// fieldName,
			// field.Type(),
			// ft.Name,
			// ft.Type,
			// typ,
		)
	}

	field.Set(child.reflectVal)

	return nil
}
