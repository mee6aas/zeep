package api

const (
	// ActivityStorage is the path of the directory that activity use.
	ActivityStorage = "/act"

	// ActivityResource is the path of the directory that contains resources activity to run.
	ActivityResource = "/act/rsc"

	// WorkflowStorage is the path of the directory shared across the workflow.
	WorkflowStorage = "/flow"

	// DockerAPIVersion is version of Docker client API that Zeep use.
	DockerAPIVersion = "1.39"

	// DefaultMongoVersion is the default version of the mongodb that Zeep use.
	DefaultMongoVersion = "4.0"

	// DefaultMongoHost is the default host name for the mongodb that Zeep use.
	DefaultMongoHost = "127.0.0.1"

	// DefaultMongoPort is the default port for the mongodb that Zeep use.
	// This value must be the same value as the default port for mongod and mongos instances.
	DefaultMongoPort = uint16(27017)
)
