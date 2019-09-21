package activity

// Descriptor describes an activity that included in the workflow.
type Descriptor struct {
	// Name of the activity
	Name string `json:"name"`

	// TODO: for parallel option
}

// Activity holds information about activity.
type Activity struct {
	Owner     string
	Name      string
	AddedDate string

	// Runtime specifies the runtime for this activity.
	Runtime string `json:"runtime"`

	// Dependencies specifies the activities included in the workflow.
	Dependencies []Descriptor `json:"dependencies"`
}
