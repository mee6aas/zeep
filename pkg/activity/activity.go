package activity

// Dep describes an activity that included in the workflow.
type Dep struct {
	// Outflow specifies where to execute the function.
	// This must be one of the "no", "optional", "always".
	//
	// - If it is "no", the function always executed in the current container
	// - If it is "optional:n", the function optionally executed in the current container.
	//   Assume executed the functions A and B executed in parallel with "optional" for outflow.
	//	 One of A and B is executed in the current container and the other one is executed in the other container.
	// - If it is "always", the function always executed in the other container.
	Outflow string `json:"outflow"`
}

// Activity holds information about activity.
type Activity struct {
	Owner     string
	Name      string
	AddedDate string

	// Runtime specifies the runtime for this activity.
	Runtime string `json:"runtime"`

	// Dependencies specifies the activities included in the workflow.
	Dependencies map[string]Dep `json:"dependencies"`

	// MaxParallelism specifies the maximum number of dependent functions that can be executed in parallel.
	// Consider that a function with `Outflow: "optional"` executed in current container is not counted.
	MaxParallelism uint `json:"maxParallelism"`
}
