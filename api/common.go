package api

const (
	// Mee6aaSDockerOrgName is the name of the M6S organization on the docker hub.
	Mee6aaSDockerOrgName = "mee6aas"

	// AgentDefaultContainerName is a default name of docker container that agent serves.
	AgentDefaultContainerName = "zeep"

	// AgentDefaultPort is a default port of the agent serves.
	AgentDefaultPort = 5122

	// AgentDefaultNetworkName is a default name of docker network that agent uses.
	AgentDefaultNetworkName = "microverse"

	// AgentTmpDirPathEnvKey is key of environment variable
	// represents the path of the directory on host that bound to temp directory on container.
	AgentTmpDirPathEnvKey = "AGENT_TEMP_DIR_HOST_PATH"

	// AgentNetworkEnvKey is key of environment variable
	// represents the name of the network that agent serves.
	AgentNetworkEnvKey = "AGENT_NETWORK"

	// AgentHostEnvKey is key of environment variable
	// represents the host name of the agent.
	AgentHostEnvKey = "AGENT_HOST"

	// AgentPortEnvKey is key of environment variable
	// represents the port of the agent.
	AgentPortEnvKey = "AGENT_PORT"

	// RuntimeDirectory is the path of the direcotry that the runtime uses.
	// This is not a host mounted storage. It should be included in the runtime image.
	RuntimeDirectory = "/runtime"

	// RuntimeResources is the path of the directory that contains the resources to run the runtime.
	RuntimeResources = "/runtime/rsc"

	// RuntimeSetup is the path of the executable that prepares the runtime to run properly.
	// This is optional. If you don't have anything to prepare before running the runtime, it can be omitted.
	RuntimeSetup = "/runtime/setup"

	// RuntimeSpawn is the path of the executable that executes the runtime.
	RuntimeSpawn = "/runtime/spawn"

	// ActivityStorage is the path of the directory that activity uses.
	// If the invocation ID is `INVOKE_ID`, the invoked activity can store data in /act/`INVOKE_ID`.
	// A special directory, where the /act/cur directory points to the directory where the currently invoked activity can store data.
	ActivityStorage = "/act"

	// ActivityResource is the path of the directory that contains resources to run activity.
	// If the activity name is `ACTIVITY_NAME`, the resource to load activity is in /act/rsc/`ACTIVITY_NAME`.
	ActivityResource = "/act/rsc"

	// ActivityManifestName is the name of the file that specifies the activity.
	ActivityManifestName = "activity.json"

	// WorkflowStorage is the path of the directory shared across the workflow.
	WorkflowStorage = "/act/flow"

	// DockerAPIVersion is version of Docker client API that Zeep uses.
	DockerAPIVersion = "1.39"
)
