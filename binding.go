package gdi

import "reflect"

// BindingKey may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type BindingKey interface{}

// BindingTag may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type BindingTag interface{}

// Binding represents an entry in the `container`. Each binding has a key and a
// corresponding value.
type Binding struct {
	key   BindingKey
	tags  map[BindingTag]struct{}
	deps  map[string]*Binding
	value interface{}
	// Todo: create scope iota
	scope     string
	opts      []Option
	populated bool
	reflectVal reflect.Value
	reflectTyp reflect.Type
}

// To Bind the key to a value.
func (b *Binding) To(v interface{}) *Binding {
	b.value = v
	return b
}

// Tag the binding.
func (b *Binding) Tag(t BindingTag) *Binding {
	b.tags[t] = struct{}{}
	return b
}

// Options TBD
func (b *Binding) Options(opts ...Option) {
	b.opts = opts
}

func (b *Binding) dependency(field string, d *Binding) {
	b.deps[field] = d
}

func (b *Binding) setReflect(val reflect.Value, typ reflect.Type) {
	b.reflectVal = val
	b.reflectTyp = typ 
}

// NewBinding return new Binding instance.
func NewBinding(key BindingKey) *Binding {
	return &Binding{
		key:  key,
		deps: make(map[string]*Binding),
	}
}
