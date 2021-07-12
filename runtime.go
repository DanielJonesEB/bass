package bass

type Runtime struct {
	Runtimes []RuntimeAssoc
}

type RuntimeAssoc struct {
	Platform Object   `json:"platform"`
	Runtime  Workload `json:"runtime"`
}

func (runtime Runtime) Run(workload Workload, cb Combiner) error {
	return nil
}

func LoadRuntime() (*Runtime, error) {
	// TODO: load runtimes.json
	return &Runtime{}, nil
}
