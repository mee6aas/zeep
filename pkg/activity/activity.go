package activity

const (
	// DefaultActivityManifestName indicate default name of manifest file for activity.
	DefaultActivityManifestName = "activity.json"
)

// Activity holds information about activity.
type Activity struct {
	Owner     string
	Name      string
	AddedDate string

	// Runtime specifies the runtime for this activity.
	Runtime string `json:"runtime"`
}
