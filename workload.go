package bass

type Workload struct {
	Args  List     `bass:"args"  optional:"true"`
	Stdin List     `bass:"stdin" optional:"true"`
	Env   Object   `bass:"env"   optional:"true"`
	From  string   `bass:"from"  optional:"true"`
	Out   Combiner `bass:"out" optional:"true"`
}
