package activity

const (
	// DefaultActivityManifestName indicate default name of manifest file for activity.
	DefaultActivityManifestName = "activity.json"
)

// Activity holds information about activity.
type Activity struct {
	// Name represents the name of this activity.
	Name string `json:"name"`

	// Runtime specifies the runtime for this activity.
	Runtime string `json:"runtime"`

	// ID represents the ID of this activity.
	// This field managed by agent.
	ID string `json:"-"`
}
