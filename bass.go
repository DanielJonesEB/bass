package bass

import (
	"errors"
	"io"
)

func New() (*Env, error) {
	env := NewEnv(ground)

	runtime, err := LoadRuntime()
	if err != nil {
		return nil, err
	}

	env.Set("run",
		Func("run", runtime.Run),
		`run a workload`,
		`A workload is an idempotent command whose result may be cached.`,
		`A workload must encapsulate all inputs which may change the result of the workload: arguments, stdin, environment variables.`,
		`A workload may specify a :platform object. A workload's platform will be used to select a runtime to run the workload it is attached to.`,
		`A runtime may be used to run a workload in an isolated or remote environment, such as a container or a cluster of worker machines.`,
	)

	return env, nil
}

func EvalReader(e *Env, r io.Reader) (Value, error) {
	reader := NewReader(r)

	var res Value
	for {
		val, err := reader.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		rdy := val.Eval(e, Identity)

		res, err = Trampoline(rdy)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func Trampoline(val Value) (Value, error) {
	for {
		var cont ReadyCont
		if err := val.Decode(&cont); err != nil {
			return val, nil
		}

		var err error
		val, err = cont.Go()
		if err != nil {
			return nil, err
		}
	}
}
