package gdi

// Option configures v using the functional options paradigm popularized by Rob Pike and Dave Cheney.
// If you're unfamiliar with this style,
// see https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html and
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis.
type Option interface {
	// Apply configuration on v.
	Apply(v interface{})
}

// OptionFunc implements Option interface.
type OptionFunc func(v interface{})

// Apply configuration on v.
func (o OptionFunc) Apply(v interface{}) {
	o(v)
}
