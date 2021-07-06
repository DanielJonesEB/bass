package bass

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Path interface {
	Value

	// DirectoryPath extends the path and returns either a DirectoryPath or a
	// FilePath.
	//
	// FilePath returns an error; it shouldn't be possible to extend a file path,
	// and this is most likely an accident.
	//
	// CommandPath also returns an error; extending a command doesn't make sense,
	// and this is most likely an accident.
	Extend(Path) (Path, error)
}

// ., ./foo/
type DirectoryPath struct {
	Path string
}

var _ Value = DirectoryPath{}

func (value DirectoryPath) String() string {
	return value.Path + "/"
}

func (value DirectoryPath) Equal(other Value) bool {
	var o DirectoryPath
	return other.Decode(&o) == nil && value.Path == o.Path
}

func (value DirectoryPath) Decode(dest interface{}) error {
	switch x := dest.(type) {
	case *Path:
		*x = value
		return nil
	case *DirectoryPath:
		*x = value
		return nil
	default:
		return DecodeError{
			Source:      value,
			Destination: dest,
		}
	}
}

// Eval returns the value.
func (value DirectoryPath) Eval(env *Env, cont Cont) ReadyCont {
	return cont.Call(value, nil)
}

var _ Path = DirectoryPath{}

func (dir DirectoryPath) Extend(ext Path) (Path, error) {
	// TODO: test
	switch p := ext.(type) {
	case DirectoryPath:
		return DirectoryPath{
			Path: dir.String() + p.Path,
		}, nil
	case FilePath:
		return FilePath{
			Path: dir.String() + p.Path,
		}, nil
	default:
		return nil, fmt.Errorf("impossible: extending path with %T", p)
	}
}

// ./foo
type FilePath struct {
	Path string
}

var _ Value = FilePath{}

func (value FilePath) String() string {
	return value.Path
}

func (value FilePath) Equal(other Value) bool {
	var o FilePath
	return other.Decode(&o) == nil && value.Path == o.Path
}

func (value FilePath) Decode(dest interface{}) error {
	switch x := dest.(type) {
	case *Path:
		*x = value
		return nil
	case *Combiner:
		*x = value
		return nil
	case *FilePath:
		*x = value
		return nil
	default:
		return DecodeError{
			Source:      value,
			Destination: dest,
		}
	}
}

// Eval returns the value.
func (value FilePath) Eval(env *Env, cont Cont) ReadyCont {
	return cont.Call(value, nil)
}

var _ Combiner = FilePath{}

func (combiner FilePath) Call(val Value, env *Env, cont Cont) ReadyCont {
	return cont.Call(combiner, nil)
}

var _ Path = FilePath{}

func (path FilePath) Extend(ext Path) (Path, error) {
	// TODO: test
	return nil, fmt.Errorf("cannot extend file path: %s", path.Path)
}

// .foo
type CommandPath struct {
	Command string
}

var _ Value = CommandPath{}

func (value CommandPath) String() string {
	// TODO: test
	return "." + value.Command
}

func (value CommandPath) Equal(other Value) bool {
	var o CommandPath
	return other.Decode(&o) == nil && value.Command == o.Command
}

func (value CommandPath) Decode(dest interface{}) error {
	switch x := dest.(type) {
	case *Path:
		*x = value
		return nil
	case *Combiner:
		*x = value
		return nil
	case *CommandPath:
		*x = value
		return nil
	default:
		return DecodeError{
			Source:      value,
			Destination: dest,
		}
	}
}

// Eval returns the value.
func (value CommandPath) Eval(env *Env, cont Cont) ReadyCont {
	// TODO: test
	return cont.Call(Applicative{value}, nil)
}

var _ Combiner = CommandPath{}

func (combiner CommandPath) Call(val Value, env *Env, cont Cont) ReadyCont {
	stdinArgs := []Value{}

	var list List
	err := val.Decode(&list)
	if err != nil {
		return cont.Call(nil, ErrBadSyntax)
	}

	kwArgs := Object{}
	var kw Keyword
	for list != (Empty{}) {
		arg := list.First()

		err = list.Rest().Decode(&list)
		if err != nil {
			return cont.Call(nil, ErrBadSyntax)
		}

		if kw != "" {
			kwArgs[kw] = arg
			kw = ""
			continue
		} else {
			var ok bool
			kw, ok = arg.(Keyword)
			if ok {
				continue
			}
		}

		stdinArgs = append(stdinArgs, arg)
	}

	var workload Workload
	if err := kwArgs.Decode(&workload); err != nil {
		return cont.Call(nil, err)
	}

	var argv []string
	if workload.Args != nil {
		err = Each(workload.Args, func(arg Value) error {
			var p Path
			if err := arg.Decode(&p); err == nil {
				argv = append(argv, p.String())
				return nil
			}

			var s string
			if err := arg.Decode(&s); err != nil {
				return err
			}

			argv = append(argv, s)

			return nil
		})
		if err != nil {
			return cont.Call(nil, err)
		}
	}

	envVars := []string{}
	if workload.Env != nil {
		// TODO: deterministic order
		for envVar, val := range workload.Env {
			var p Path
			if err := val.Decode(&p); err == nil {
				envVars = append(envVars, fmt.Sprintf("%s=%s", string(envVar), p))
				continue
			}

			var s string
			if err := val.Decode(&s); err != nil {
				return cont.Call(nil, err)
			}

			envVars = append(envVars, fmt.Sprintf("%s=%s", string(envVar), s))
		}
	}

	cmd := exec.Command(combiner.Command, argv...)
	cmd.Stderr = Stderr
	cmd.Env = append(os.Environ(), envVars...)
	cmd.Dir = workload.From

	in, err := cmd.StdinPipe()
	if err != nil {
		return cont.Call(nil, err)
	}

	var outStream Value
	if workload.Out != nil {
		out, err := cmd.StdoutPipe()
		if err != nil {
			return cont.Call(nil, err)
		}

		outStream = &Source{
			NewJSONSource(combiner.Command, out),
		}
	} else {
		cmd.Stdout = Stderr
	}

	err = cmd.Start()
	if err != nil {
		return cont.Call(nil, err)
	}

	enc := json.NewEncoder(in)
	for _, arg := range stdinArgs {
		err := enc.Encode(arg)
		if err != nil {
			return cont.Call(nil, err)
		}
	}

	err = in.Close()
	if err != nil {
		return cont.Call(nil, err)
	}

	if outStream != nil {
		return workload.Out.Call(NewList(outStream), env, Continue(func(res Value) Value {
			err = cmd.Wait()
			if err != nil {
				return cont.Call(nil, err)
			}

			return cont.Call(res, nil)
		}))
	}

	err = cmd.Wait()
	if err != nil {
		return cont.Call(nil, err)
	}

	return cont.Call(Null{}, nil)
}

type objectStream struct {
	closer io.Closer
	dec    *json.Decoder
}

func (stream *objectStream) Next(ctx context.Context) (Value, error) {
	var val interface{}
	err := stream.dec.Decode(&val)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, ErrEndOfSource
		}

		return nil, err
	}

	return ValueOf(val)
}

func (stream *objectStream) Close(context.Context) error {
	return stream.closer.Close()
}

var _ Path = CommandPath{}

func (path CommandPath) Extend(ext Path) (Path, error) {
	// TODO: test
	return nil, fmt.Errorf("cannot extend command path: %s", path.Command)
}

type ExtendPath struct {
	Parent Value
	Child  Path
}

func (value ExtendPath) String() string {
	// TODO: test
	return fmt.Sprintf("%s%s", value.Parent, value.Child)
}

func (value ExtendPath) Equal(other Value) bool {
	// TODO: test
	var o ExtendPath
	if err := other.Decode(&o); err != nil {
		return false
	}

	return value.Parent.Equal(o.Parent) && value.Child.Equal(o.Child)
}

func (value ExtendPath) Decode(dest interface{}) error {
	switch x := dest.(type) {
	case *ExtendPath:
		*x = value
		return nil
	default:
		return DecodeError{
			Source:      value,
			Destination: dest,
		}
	}
}

// Eval returns the value.
func (value ExtendPath) Eval(env *Env, cont Cont) ReadyCont {
	// TODO: test
	return value.Parent.Eval(env, Continue(func(parent Value) Value {
		var path Path
		if err := parent.Decode(&path); err != nil {
			return cont.Call(nil, err)
		}

		return cont.Call(path.Extend(value.Child))
	}))
}
