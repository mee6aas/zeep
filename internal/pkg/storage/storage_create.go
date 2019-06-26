package storage

// CreateConfig holds the configuration for the Create operation.
type CreateConfig struct {
	Size uint32 // Size limit of the storage to create
}

// CreateCreatedBody OK response to Create operation.
type CreateCreatedBody struct {
	Path string // Path of the storage
}
