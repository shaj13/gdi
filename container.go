package gdi

import (
	"fmt"
	"reflect"
)

// Container TBD
type Container struct {
	parent      *Container
	name        string
	registry    map[BindingKey]*Binding
	unpopulated map[BindingKey]*Binding
}

// Bind create a new binding and add it to the  conatiner.
func (c *Container) Bind(key BindingKey) *Binding {
	b := NewBinding(key)
	c.unpopulated[key] = b
	return b
}

// Populate TBD.
func (c *Container) Populate() error {
	// if c.parent != nil {
	// 	if err := c.parent.Populate(); err != nil {
	// 		return err
	// 	}
	// }
	for k := range c.unpopulated {
		_, err := c.resolve(k)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Container) resolve(key BindingKey) (*Binding, error) {
	b, ok := c.unpopulated[key]

	if !ok {
		return nil, fmt.Errorf("gdi: The key '%v' is not bound to any value in container", key)
	}

	value := reflect.ValueOf(b.value)
	typ := reflect.TypeOf(b.value)

	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		for i := 0; i < value.Elem().NumField(); i++ {
			key, ok := typ.Elem().Field(i).Tag.Lookup("inject")

			if !ok {
				continue
			}

			fb, err := c.resolve(BindingKey(key))
			if err != nil {
				return nil, err
			}

			fv := value.Elem().Field(i)
			ft := typ.Elem().Field(i)

			if !fv.CanSet() {
				return nil, fmt.Errorf(
					"gdi: inject on unexported field %s in object of type %s",
					ft.Name,
					typ,
				)
			}

			fbt := reflect.TypeOf(fb.value)
			fbv := reflect.ValueOf(fb.value)



			if ft.Type.Kind() == reflect.Interface && fbt.Implements(ft.Type) {
				reflect.Indirect(fv).Set(fbv)
				continue
			}

			if fv.Kind() == reflect.Struct && fbv.Kind() == reflect.Ptr {
				fbv = fbv.Elem()
				fbt = fbv.Type()
			}
			
			if !ft.Type.AssignableTo(fbt) {
				return nil, fmt.Errorf(
					"gdi: inject %s of type %s is not assignable to field %s of type %s in object of type %s",
					key,
					fbt,
					ft.Name,
					ft.Type,
					typ,
				)
			}

			fv.Set(fbv)
		}
	}

	c.registry[b.key] = b
	return b, nil
}

// func isStructPtr(t reflect.Type) bool {
// 	return
// }

func NewContainer() *Container {
	return &Container{
		name:        "sample",
		registry:    make(map[BindingKey]*Binding),
		unpopulated: make(map[BindingKey]*Binding),
	}
}