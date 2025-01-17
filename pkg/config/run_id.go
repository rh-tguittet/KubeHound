package config

import (
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// RunID represents a unique ID for each KubeHound run.
type RunID struct {
	val ulid.ULID
}

// NewRunID creates a new RunID instance.
func NewRunID() *RunID {
	return &RunID{
		val: ulid.Make(),
	}
}

// String returns the string representation of the run id.
// NOTE: this is lowercased to ensure consistency with Datadog (where tags are automatically lower cased)
func (r RunID) String() string {
	return strings.ToLower(r.val.String())
}

// Timestamp returns the timestamp embedded within the run id.
func (r RunID) Timestamp() time.Time {
	return time.UnixMilli(int64(r.val.Time()))
}
