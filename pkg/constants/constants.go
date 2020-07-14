package constants

const (
	Version               = "v0.0.1"
	ApiServer             = "https://staroid.com/api"
	EnvStaroidAccessToken = "STAROID_ACCESS_TOKEN"
	EnvStaroidApiServer   = "STAROID_API_SERVER"
	TunnelServicePort     = 57682
	KubeproxyPort         = 57683

	K8S_LABEL_KEY_RESOURCE_SYSTEM         = "resource.staroid.com/system"
	K8S_LABEL_VALUE_RESOURCE_SYSTEM_SHELL = "shell"

	StatusPollingIntervalSec = 5
	NsStartTimeoutSec        = 10 * 60
)
