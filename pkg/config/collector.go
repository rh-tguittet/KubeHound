package config

const (
	CollectorTypeFile         = "file-collector"
	CollectorTypeK8sAPI       = "live-k8s-api-collector"
	CollectorTypeOpenShiftAPI = "live-openshift-api-collector"
)

const (
	DefaultK8sAPIPageSize           int64 = 500
	DefaultK8sAPIPageBufferSize     int32 = 10
	DefaultK8sAPIRateLimitPerSecond int   = 100
)

// CollectorConfig configures collector specific parameters.
type CollectorConfig struct {
	Type          string                       `mapstructure:"type"`           // Collector type
	File          *FileCollectorConfig         `mapstructure:"file"`           // File collector specific configuration
	Live          *K8SAPICollectorConfig       `mapstructure:"live"`           // K8S collector specific configuration
	LiveOpenShift *OpenShiftAPICollectorConfig `mapstructure:"live-openshift"` // OpenShift collector specific configuration
}

// K8SAPICollectorConfig configures the K8sAPI collector.
type K8SAPICollectorConfig struct {
	PageSize           int64 `mapstructure:"page_size"`             // Number of entry being retrieving by each call on the API (same for all Kubernetes entry types)
	PageBufferSize     int32 `mapstructure:"page_buffer_size"`      // Number of pages to buffer
	RateLimitPerSecond int   `mapstructure:"rate_limit_per_second"` // Rate limiting per second across all calls (same for all kubernetes entry types) against the Kubernetes API
}

// FileCollectorConfig configures the file collector.
type FileCollectorConfig struct {
	ClusterName string `mapstructure:"cluster"`   // Target cluster (must be specified in config as not present in JSON files)
	Directory   string `mapstructure:"directory"` // Base directory holding the K8s data JSON files
}

// OpenShiftAPICollectorConfig builds upon K8SAPICollectorConfig
type OpenShiftAPICollectorConfig struct {
	K8SAPICollectorConfig
}
