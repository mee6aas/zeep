package api

const (
	// Mee6aaSDockerOrgName is the name of the docker hub organization Mee6aaS use.
	Mee6aaSDockerOrgName = "mee6aas"

	// NetworkName is a name of docker network that agent use.
	NetworkName = "m6s"

	// ActivityStorage is the path of the directory that activity use.
	// If the invocation ID is `INVOKE_ID`, the invoked activity can store data in /act/`INVOKE_ID`.
	// A special directory, where the /act/cur directory points to the directory where the currently invoked activity can store data.
	ActivityStorage = "/act"

	// ActivityResource is the path of the directory that contains resources to run activity.
	// If the activity name is `ACTIVITY_NAME`, the resource to load activity is in /act/rsc/`ACTIVITY_NAME`.
	ActivityResource = "/act/rsc"

	// WorkflowStorage is the path of the directory shared across the workflow.
	WorkflowStorage = "/act/flow"

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
