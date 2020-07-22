package gdi

import (
	"fmt"
	"reflect"
)

type walker map[BindingKey]*Binding

func (w walker) walk(b *Binding) error {
	if b.populated {
		return nil
	}

	_, err := w.walkAndResolve(b.key)
	return err
}

func (w walker) walkAndResolve(key BindingKey) (*Binding, error) {
	binding, ok := w[key]

	if !ok {
		return nil, fmt.Errorf("gdi: The key '%v' is not bound to any value in container", key)
	}

	val := reflect.ValueOf(binding.value)
	typ := reflect.TypeOf(binding.value)
	binding.setReflect(val, typ)

	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		for i := 0; i < val.Elem().NumField(); i++ {
			field := typ.Elem().Field(i)
			bindKey, ok := field.Tag.Lookup("inject")

			if !ok {
				continue
			}

			fieldBinding, err := w.walkAndResolve(bindKey)

			if err != nil {
				return nil, err
			}

			err = transient(binding, fieldBinding, field.Name)
			if err != nil {
				return nil, err
			}

			binding.dependency(field.Name, fieldBinding)
		}
	}

	return binding, nil
}
