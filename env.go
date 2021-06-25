package bass

// Env contains bindings from symbols to values, and parent environments to
// delegate to during symbol lookup.
type Env struct {
	Bindings map[Symbol]Value
	Parents  []*Env
}

// Assert that Env is a Value.
var _ Value = (*Env)(nil)

// NewEnv constructs an Env with empty bindings and the given parents.
func NewEnv(parents ...*Env) *Env {
	return &Env{
		Bindings: map[Symbol]Value{},
		Parents:  parents,
	}
}

func (value *Env) String() string {
	return "<env>"
}

func (value *Env) Decode(dest interface{}) error {
	switch x := dest.(type) {
	case **Env:
		*x = value
		return nil
	}

	return DecodeError{
		Source:      value,
		Destination: dest,
	}
}

// Eval returns the value.
func (value *Env) Eval(env *Env) (Value, error) {
	return value, nil
}

// Set assigns the value in the local bindings.
func (env *Env) Set(binding Symbol, value Value) {
	env.Bindings[binding] = value
}

// Get fetches the given binding.
//
// If a value is set in the local bindings, it is returned.
//
// If not, the parent environments are queried in order.
//
// If no value is found, false is returned.
func (env *Env) Get(binding Symbol) (Value, bool) {
	val, found := env.Bindings[binding]
	if found {
		return val, found
	}

	for _, parent := range env.Parents {
		val, found = parent.Get(binding)
		if found {
			return val, found
		}
	}

	return nil, false
}
