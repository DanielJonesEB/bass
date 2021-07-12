package bass

// A workload is a black-box object containing fields interpreted by its
// runtime.
type Workload Object

// HostWorkload is a command to run natively on the host machine.
type HostWorkload struct {
	// Path is either a string or a path value specifying the command or file to run.
	Path Value `bass:"path"`

	// Args is a list of string or path values, including argv[0] as the program
	// being run.
	Args List `bass:"args" optional:"true"`

	// Stdin is a fixed list of values to write as a JSON stream on stdin.
	//
	// This is distinct from a stream interface; it is a finite part of the
	// request so that it may be used to form a cache key.
	Stdin List `bass:"stdin" optional:"true"`

	// Env is a map of environment variables to set for the workload.
	Env Object `bass:"env" optional:"true"`

	// From is a string or path value specifying the working directory the process should be run within.
	Dir Value `bass:"dir" optional:"true"`
}
