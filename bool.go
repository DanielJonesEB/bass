package bass

import "context"

type Bool bool

func (value Bool) String() string {
	if bool(value) {
		return "true"
	} else {
		return "false"
	}
}

func (value Bool) Equal(other Value) bool {
	var o Bool
	return other.Decode(&o) == nil && value == o
}

func (value Bool) Decode(dest interface{}) error {
	switch x := dest.(type) {
	case *Bool:
		*x = value
		return nil
	case *Bindable:
		*x = value
		return nil
	case *Value:
		*x = value
		return nil
	case *bool:
		*x = bool(value)
		return nil
	default:
		return DecodeError{
			Source:      value,
			Destination: dest,
		}
	}
}

// Eval returns the value.
func (value Bool) Eval(ctx context.Context, env *Env, cont Cont) ReadyCont {
	return cont.Call(value, nil)
}

var _ Bindable = Bool(false)

func (binding Bool) Bind(env *Env, val Value) error {
	return BindConst(binding, val)
}
