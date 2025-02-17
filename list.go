package bass

import "fmt"

type List interface {
	Value
	Bindable

	First() Value
	Rest() Value
}

func NewList(vals ...Value) List {
	var list List = Empty{}
	for i := len(vals) - 1; i >= 0; i-- {
		list = Pair{
			A: vals[i],
			D: list,
		}
	}

	return list
}

func Each(list List, cb func(Value) error) error {
	for list != (Empty{}) {
		err := cb(list.First())
		if err != nil {
			return err
		}

		err = list.Rest().Decode(&list)
		if err != nil {
			// TODO: better error
			return fmt.Errorf("each: %w", err)
		}
	}

	return nil
}

func ToSlice(list List) ([]Value, error) {
	var vals []Value
	err := Each(list, func(v Value) error {
		vals = append(vals, v)
		return nil
	})
	if err != nil {
		// malformed list
		return nil, fmt.Errorf("to slice: %w", err)
	}

	return vals, nil
}

func IsList(val Value) bool {
	var empty Empty
	err := val.Decode(&empty)
	if err == nil {
		return true
	}

	var list List
	err = val.Decode(&list)
	if err != nil {
		return false
	}

	return IsList(list.Rest())
}

func BindList(binding List, env *Env, value Value) error {
	mismatch := BindMismatchError{
		Need: binding,
		Have: value,
	}

	var e Empty
	if err := value.Decode(&e); err == nil {
		// empty value given for list
		return mismatch
	}

	var v List
	if err := value.Decode(&v); err != nil {
		// non-list value
		return mismatch
	}

	var f Bindable
	if err := binding.First().Decode(&f); err != nil {
		return CannotBindError{binding.First()}
	}

	if err := f.Bind(env, v.First()); err != nil {
		return err
	}

	var r Bindable
	if err := binding.Rest().Decode(&r); err != nil {
		return CannotBindError{binding.Rest()}
	}

	return r.Bind(env, v.Rest())
}
