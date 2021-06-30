package bass

type Symbol string

func (value Symbol) String() string {
	return string(value)
}

func (value Symbol) Equal(other Value) bool {
	var o Symbol
	return other.Decode(&o) == nil && value == o
}

func (value Symbol) Decode(dest interface{}) error {
	switch x := dest.(type) {
	case *Symbol:
		*x = value
		return nil
	case *string:
		*x = string(value)
		return nil
	default:
		return DecodeError{
			Source:      value,
			Destination: dest,
		}
	}
}

// Eval returns the value.
func (value Symbol) Eval(env *Env) (Value, error) {
	res, found := env.Get(value)
	if !found {
		return nil, UnboundError{value}
	}

	return res, nil
}
